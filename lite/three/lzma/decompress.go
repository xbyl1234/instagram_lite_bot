package lzma

import (
	"io"
	"unsafe"
)

func (r *Reader1) decompress(needBytesCount uint32) (err error) {
	s := r.s
	rCode := r.rangeDec.Code
	rRange := r.rangeDec.Range

	for r.outWindow.pending < needBytesCount {
		if s.unpackSizeDefined && s.bytesLeft == 0 {
			if rCode == 0 {
				err = io.EOF

				break
			}
		}

		s.posState = r.outWindow.pos & s.posMask
		state2 := (s.state << kNumPosBitsMax) + s.posState

		{ // r.rangeDec.DecodeBit(&s.isMatch[state2])
			v := (*prob)(unsafe.Pointer(uintptr(unsafe.Pointer(&s.isMatch[0])) + uintptr(state2)*unsafe.Sizeof(prob(0))))
			bound := (rRange >> kNumBitModelTotalBits) * uint32(*v)

			if rCode < bound {
				*v += ((1 << kNumBitModelTotalBits) - *v) >> kNumMoveBits
				rRange = bound

				// Normalize
				if rRange < kTopValue {
					b, err := r.rangeDec.inStream.ReadByte()
					if err != nil {
						return err
					}

					rRange <<= 8
					rCode = (rCode << 8) | uint32(b)
				}

				{ // literal
					if s.unpackSizeDefined && s.bytesLeft == 0 {
						return ErrResultError
					}

					{ // DecodeLiteral
						prevByte := uint32(0)
						if !r.outWindow.IsEmpty() {
							prevByte = uint32(r.outWindow.GetByte(1))
						}

						symbol := uint32(1)
						litState := ((r.outWindow.pos & ((1 << s.lp) - 1)) << s.lc) + (prevByte >> (8 - s.lc))
						probsPtr := uintptr(unsafe.Pointer(&s.litProbs[0])) + uintptr(0x300*litState)*unsafe.Sizeof(prob(0))

						if s.state >= 7 {
							matchByte := r.outWindow.GetByte(s.rep0 + 1)

							for symbol < 0x100 {
								matchBit := uint32((matchByte >> 7) & 1)
								matchByte <<= 1
								probPtr := (*prob)(unsafe.Pointer(probsPtr + uintptr(((1+matchBit)<<8)+symbol)*unsafe.Sizeof(prob(0))))

								{ // rc.DecodeBit
									bound := (rRange >> kNumBitModelTotalBits) * uint32(*probPtr)

									if rCode < bound {
										*probPtr += ((1 << kNumBitModelTotalBits) - *probPtr) >> kNumMoveBits
										rRange = bound
										symbol <<= 1

										if matchBit != 0 {
											break
										}

										// Normalize
										if rRange < kTopValue {
											b, err := r.rangeDec.inStream.ReadByte()
											if err != nil {
												return err
											}

											rRange <<= 8
											rCode = (rCode << 8) | uint32(b)
										}
									} else {
										*probPtr -= *probPtr >> kNumMoveBits
										rCode -= bound
										rRange -= bound
										symbol = (symbol << 1) | 1

										if matchBit != 1 {
											break
										}

										// Normalize
										if rRange < kTopValue {
											b, err := r.rangeDec.inStream.ReadByte()
											if err != nil {
												return err
											}

											rRange <<= 8
											rCode = (rCode << 8) | uint32(b)
										}
									}
								}
							}
						}

						// Normalize
						if rRange < kTopValue {
							b, err := r.rangeDec.inStream.ReadByte()
							if err != nil {
								return err
							}

							rRange <<= 8
							rCode = (rCode << 8) | uint32(b)
						}

						for symbol < 0x100 {
							probPtr := (*prob)(unsafe.Pointer(probsPtr + uintptr(symbol)*unsafe.Sizeof(prob(0))))
							{ // rc.DecodeBit
								bound := (rRange >> kNumBitModelTotalBits) * uint32(*probPtr)

								if rCode < bound {
									*probPtr += ((1 << kNumBitModelTotalBits) - *probPtr) >> kNumMoveBits
									rRange = bound
									symbol <<= 1

									// Normalize
									if rRange < kTopValue {
										b, err := r.rangeDec.inStream.ReadByte()
										if err != nil {
											return err
										}

										rRange <<= 8
										rCode = (rCode << 8) | uint32(b)
									}
								} else {
									*probPtr -= *probPtr >> kNumMoveBits
									rCode -= bound
									rRange -= bound
									symbol = (symbol << 1) | 1

									// Normalize
									if rRange < kTopValue {
										b, err := r.rangeDec.inStream.ReadByte()
										if err != nil {
											return err
										}

										rRange <<= 8
										rCode = (rCode << 8) | uint32(b)
									}
								}
							}
						}

						r.outWindow.PutByte(byte(symbol - 0x100))
					}

					s.state = stateUpdateLiteral(s.state)
					s.bytesLeft--

					continue
				}
			} else {
				*v -= *v >> kNumMoveBits
				rCode -= bound
				rRange -= bound

				// Normalize
				if rRange < kTopValue {
					b, err := r.rangeDec.inStream.ReadByte()
					if err != nil {
						return err
					}

					rRange <<= 8
					rCode = (rCode << 8) | uint32(b)
				}

				{ // match
					length := uint32(0)

					{ // r.rangeDec.DecodeBit(&s.isRep[s.state])
						v := (*prob)(unsafe.Pointer(uintptr(unsafe.Pointer(&s.isRep[0])) + uintptr(s.state)*unsafe.Sizeof(prob(0))))
						bound := (rRange >> kNumBitModelTotalBits) * uint32(*v)

						if rCode < bound {
							*v += ((1 << kNumBitModelTotalBits) - *v) >> kNumMoveBits
							rRange = bound

							// Normalize
							if rRange < kTopValue {
								b, err := r.rangeDec.inStream.ReadByte()
								if err != nil {
									return err
								}

								rRange <<= 8
								rCode = (rCode << 8) | uint32(b)
							}

							{ // simple match
								s.rep3, s.rep2, s.rep1 = s.rep2, s.rep1, s.rep0

								{ // lenDecoder.Decode
									{ // r.rangeDec.DecodeBit(&s.lenDecoder.choice)
										bound := (rRange >> kNumBitModelTotalBits) * uint32(s.lenDecoderChoice)

										if rCode < bound {
											s.lenDecoderChoice += ((1 << kNumBitModelTotalBits) - s.lenDecoderChoice) >> kNumMoveBits
											rRange = bound

											// Normalize
											if rRange < kTopValue {
												b, err := r.rangeDec.inStream.ReadByte()
												if err != nil {
													return err
												}

												rRange <<= 8
												rCode = (rCode << 8) | uint32(b)
											}

											{ // s.lenDecoder.lowCoder[s.posState].Decode
												m := uint32(1)

												probsPtr := uintptr(unsafe.Pointer(&s.lenDecoderLowCoder[s.posState][0]))

												for i := 0; i < lenLowCoderNumBits; i++ {
													probPtr := (*prob)(unsafe.Pointer(probsPtr + uintptr(m)*unsafe.Sizeof(prob(0))))
													{ // rc.DecodeBit
														bound := (rRange >> kNumBitModelTotalBits) * uint32(*probPtr)

														if rCode < bound {
															*probPtr += ((1 << kNumBitModelTotalBits) - *probPtr) >> kNumMoveBits
															rRange = bound
															m <<= 1

															// Normalize
															if rRange < kTopValue {
																b, err := r.rangeDec.inStream.ReadByte()
																if err != nil {
																	return err
																}

																rRange <<= 8
																rCode = (rCode << 8) | uint32(b)
															}
														} else {
															*probPtr -= *probPtr >> kNumMoveBits
															rCode -= bound
															rRange -= bound
															m = (m << 1) | 1

															// Normalize
															if rRange < kTopValue {
																b, err := r.rangeDec.inStream.ReadByte()
																if err != nil {
																	return err
																}

																rRange <<= 8
																rCode = (rCode << 8) | uint32(b)
															}
														}
													}
												}

												length = m - (uint32(1) << lenLowCoderNumBits)
											}
										} else {
											s.lenDecoderChoice -= s.lenDecoderChoice >> kNumMoveBits
											rCode -= bound
											rRange -= bound

											// Normalize
											if rRange < kTopValue {
												b, err := r.rangeDec.inStream.ReadByte()
												if err != nil {
													return err
												}

												rRange <<= 8
												rCode = (rCode << 8) | uint32(b)
											}

											{ // r.rangeDec.DecodeBit(&s.lenDecoder.choice2)
												bound := (rRange >> kNumBitModelTotalBits) * uint32(s.lenDecoderChoice2)

												if rCode < bound {
													s.lenDecoderChoice2 += ((1 << kNumBitModelTotalBits) - s.lenDecoderChoice2) >> kNumMoveBits
													rRange = bound

													// Normalize
													if rRange < kTopValue {
														b, err := r.rangeDec.inStream.ReadByte()
														if err != nil {
															return err
														}

														rRange <<= 8
														rCode = (rCode << 8) | uint32(b)
													}

													{ // s.lenDecoder.midCoder[s.posState].Decode
														m := uint32(1)

														probsPtr := uintptr(unsafe.Pointer(&s.lenDecoderMidCoder[s.posState][0]))

														for i := 0; i < lenMidCoderNumBits; i++ {
															probPtr := (*prob)(unsafe.Pointer(probsPtr + uintptr(m)*unsafe.Sizeof(prob(0))))
															{ // rc.DecodeBit
																bound := (rRange >> kNumBitModelTotalBits) * uint32(*probPtr)

																if rCode < bound {
																	*probPtr += ((1 << kNumBitModelTotalBits) - *probPtr) >> kNumMoveBits
																	rRange = bound
																	m <<= 1

																	// Normalize
																	if rRange < kTopValue {
																		b, err := r.rangeDec.inStream.ReadByte()
																		if err != nil {
																			return err
																		}

																		rRange <<= 8
																		rCode = (rCode << 8) | uint32(b)
																	}
																} else {
																	*probPtr -= *probPtr >> kNumMoveBits
																	rCode -= bound
																	rRange -= bound
																	m = (m << 1) | 1

																	// Normalize
																	if rRange < kTopValue {
																		b, err := r.rangeDec.inStream.ReadByte()
																		if err != nil {
																			return err
																		}

																		rRange <<= 8
																		rCode = (rCode << 8) | uint32(b)
																	}
																}
															}
														}

														length = 8 + m - (uint32(1) << lenMidCoderNumBits)
													}
												} else {
													s.lenDecoderChoice2 -= s.lenDecoderChoice2 >> kNumMoveBits
													rCode -= bound
													rRange -= bound

													// Normalize
													if rRange < kTopValue {
														b, err := r.rangeDec.inStream.ReadByte()
														if err != nil {
															return err
														}

														rRange <<= 8
														rCode = (rCode << 8) | uint32(b)
													}

													{ // s.lenDecoder.highCoder.Decode
														m := uint32(1)

														probsPtr := uintptr(unsafe.Pointer(&s.lenDecoderHighCoder[0]))

														for i := 0; i < lenHighCoderNumBits; i++ {
															probPtr := (*prob)(unsafe.Pointer(probsPtr + uintptr(m)*unsafe.Sizeof(prob(0))))
															{ // rc.DecodeBit
																bound := (rRange >> kNumBitModelTotalBits) * uint32(*probPtr)

																if rCode < bound {
																	*probPtr += ((1 << kNumBitModelTotalBits) - *probPtr) >> kNumMoveBits
																	rRange = bound
																	m <<= 1

																	// Normalize
																	if rRange < kTopValue {
																		b, err := r.rangeDec.inStream.ReadByte()
																		if err != nil {
																			return err
																		}

																		rRange <<= 8
																		rCode = (rCode << 8) | uint32(b)
																	}
																} else {
																	*probPtr -= *probPtr >> kNumMoveBits
																	rCode -= bound
																	rRange -= bound
																	m = (m << 1) | 1

																	// Normalize
																	if rRange < kTopValue {
																		b, err := r.rangeDec.inStream.ReadByte()
																		if err != nil {
																			return err
																		}

																		rRange <<= 8
																		rCode = (rCode << 8) | uint32(b)
																	}
																}
															}
														}

														length = 16 + m - (uint32(1) << lenHighCoderNumBits)
													}
												}
											}
										}
									}
								}

								s.state = stateUpdateMatch(s.state)

								{ // DecodeDistance
									lenState := length
									if lenState > (kNumLenToPosStates - 1) {
										lenState = kNumLenToPosStates - 1
									}

									var posSlot uint32

									{ // s.posSlotDecoder[lenState].Decode
										m := uint32(1)

										probsPtr := uintptr(unsafe.Pointer(&s.posSlotDecoderProbs[lenState][0]))

										for i := 0; i < posSlotDecoderNumBits; i++ {
											probPtr := (*prob)(unsafe.Pointer(probsPtr + uintptr(m)*unsafe.Sizeof(prob(0))))
											{ // rc.DecodeBit
												bound := (rRange >> kNumBitModelTotalBits) * uint32(*probPtr)

												if rCode < bound {
													*probPtr += ((1 << kNumBitModelTotalBits) - *probPtr) >> kNumMoveBits
													rRange = bound
													m <<= 1

													// Normalize
													if rRange < kTopValue {
														b, err := r.rangeDec.inStream.ReadByte()
														if err != nil {
															return err
														}

														rRange <<= 8
														rCode = (rCode << 8) | uint32(b)
													}
												} else {
													*probPtr -= *probPtr >> kNumMoveBits
													rCode -= bound
													rRange -= bound
													m = (m << 1) | 1

													// Normalize
													if rRange < kTopValue {
														b, err := r.rangeDec.inStream.ReadByte()
														if err != nil {
															return err
														}

														rRange <<= 8
														rCode = (rCode << 8) | uint32(b)
													}
												}
											}
										}

										posSlot = m - (uint32(1) << posSlotDecoderNumBits)
									}

									if posSlot < 4 {
										s.rep0 = posSlot
									} else {
										numDirectBits := (posSlot >> 1) - 1
										dist := (2 | (posSlot & 1)) << numDirectBits

										if posSlot < kEndPosModelIndex {
											{ // BitTreeReverseDecode
												probsPtr := uintptr(unsafe.Pointer(&s.posDecoders[0])) + uintptr(dist-posSlot)*unsafe.Sizeof(prob(0))

												m := uint32(1)
												symbol := uint32(0)

												for i := uint32(0); i < numDirectBits; i++ {
													probPtr := (*prob)(unsafe.Pointer(probsPtr + uintptr(m)*unsafe.Sizeof(prob(0))))
													{ // rc.DecodeBit
														bound := (rRange >> kNumBitModelTotalBits) * uint32(*probPtr)

														if rCode < bound {
															*probPtr += ((1 << kNumBitModelTotalBits) - *probPtr) >> kNumMoveBits
															rRange = bound
															m <<= 1
															symbol |= 0 << i

															// Normalize
															if rRange < kTopValue {
																b, err := r.rangeDec.inStream.ReadByte()
																if err != nil {
																	return err
																}

																rRange <<= 8
																rCode = (rCode << 8) | uint32(b)
															}
														} else {
															*probPtr -= *probPtr >> kNumMoveBits
															rCode -= bound
															rRange -= bound
															m = (m << 1) | 1
															symbol |= 1 << i

															// Normalize
															if rRange < kTopValue {
																b, err := r.rangeDec.inStream.ReadByte()
																if err != nil {
																	return err
																}

																rRange <<= 8
																rCode = (rCode << 8) | uint32(b)
															}
														}
													}
												}

												dist += symbol
												s.rep0 = dist
											}
										} else {
											var res uint32
											{ // DecodeDirectBits
												numBits := numDirectBits - kNumAlignBits

												for ; numBits > 0; numBits-- {
													rRange >>= 1
													rCode -= rRange
													t := 0 - (rCode >> 31)
													rCode += rRange & t

													if rCode == rRange {
														r.rangeDec.Corrupted = true
													}

													res <<= 1
													res += t + 1

													// Normalize
													if rRange < kTopValue {
														b, err := r.rangeDec.inStream.ReadByte()
														if err != nil {
															return err
														}

														rRange <<= 8
														rCode = (rCode << 8) | uint32(b)
													}
												}
											}
											dist += res << kNumAlignBits

											symbol := uint32(0)
											{ // BitTreeReverseDecode
												probsPtr := uintptr(unsafe.Pointer(&s.alignDecoderProbs[0]))

												m := uint32(1)

												for i := 0; i < kNumAlignBits; i++ {
													probPtr := (*prob)(unsafe.Pointer(probsPtr + uintptr(m)*unsafe.Sizeof(prob(0))))
													{ // rc.DecodeBit
														bound := (rRange >> kNumBitModelTotalBits) * uint32(*probPtr)

														if rCode < bound {
															*probPtr += ((1 << kNumBitModelTotalBits) - *probPtr) >> kNumMoveBits
															rRange = bound
															m <<= 1
															symbol |= 0 << i

															// Normalize
															if rRange < kTopValue {
																b, err := r.rangeDec.inStream.ReadByte()
																if err != nil {
																	return err
																}

																rRange <<= 8
																rCode = (rCode << 8) | uint32(b)
															}
														} else {
															*probPtr -= *probPtr >> kNumMoveBits
															rCode -= bound
															rRange -= bound
															m = (m << 1) | 1
															symbol |= 1 << i

															// Normalize
															if rRange < kTopValue {
																b, err := r.rangeDec.inStream.ReadByte()
																if err != nil {
																	return err
																}

																rRange <<= 8
																rCode = (rCode << 8) | uint32(b)
															}
														}
													}
												}
											}

											dist += symbol
											s.rep0 = dist
										}
									}
								}

								if s.rep0 == 0xFFFFFFFF {
									if rCode == 0 {
										if s.unpackSizeDefined && s.bytesLeft > 0 {
											return ErrResultError
										}

										err = io.EOF

										break
									} else {
										return ErrResultError
									}
								}

								if s.unpackSizeDefined && s.bytesLeft == 0 {
									return ErrResultError
								}

								if s.rep0 >= r.outWindow.size || !r.outWindow.CheckDistance(s.rep0) {
									return ErrResultError
								}
							}

							length += kMatchMinLen
							if s.unpackSizeDefined && uint32(s.bytesLeft) < length {
								length = uint32(s.bytesLeft)
								r.outWindow.CopyMatch(s.rep0+1, length)
								s.bytesLeft -= uint64(length)

								return ErrResultError
							} else {
								r.outWindow.CopyMatch(s.rep0+1, length)
								s.bytesLeft -= uint64(length)

								continue
							}
						} else {
							*v -= *v >> kNumMoveBits
							rCode -= bound
							rRange -= bound

							// Normalize
							if rRange < kTopValue {
								b, err := r.rangeDec.inStream.ReadByte()
								if err != nil {
									return err
								}

								rRange <<= 8
								rCode = (rCode << 8) | uint32(b)
							}

							{ // rep match
								if s.unpackSizeDefined && s.bytesLeft == 0 {
									return ErrResultError
								}

								if r.outWindow.IsEmpty() {
									return ErrResultError
								}

								{ // r.rangeDec.DecodeBit(&s.isRepG0[s.state])
									v := (*prob)(unsafe.Pointer(uintptr(unsafe.Pointer(&s.isRepG0[0])) + uintptr(s.state)*unsafe.Sizeof(prob(0))))
									bound := (rRange >> kNumBitModelTotalBits) * uint32(*v)

									if rCode < bound {
										*v += ((1 << kNumBitModelTotalBits) - *v) >> kNumMoveBits
										rRange = bound

										// Normalize
										if rRange < kTopValue {
											b, err := r.rangeDec.inStream.ReadByte()
											if err != nil {
												return err
											}

											rRange <<= 8
											rCode = (rCode << 8) | uint32(b)
										}

										{ // short rep match
											{ // r.rangeDec.DecodeBit(&s.isRep0Long[state2])
												v := (*prob)(unsafe.Pointer(uintptr(unsafe.Pointer(&s.isRep0Long[0])) + uintptr(state2)*unsafe.Sizeof(prob(0))))
												bound := (rRange >> kNumBitModelTotalBits) * uint32(*v)

												if rCode < bound {
													*v += ((1 << kNumBitModelTotalBits) - *v) >> kNumMoveBits
													rRange = bound

													// Normalize
													if rRange < kTopValue {
														b, err := r.rangeDec.inStream.ReadByte()
														if err != nil {
															return err
														}

														rRange <<= 8
														rCode = (rCode << 8) | uint32(b)
													}

													s.state = stateUpdateShortRep(s.state)
													r.outWindow.PutByte(r.outWindow.GetByte(s.rep0 + 1))
													s.bytesLeft--

													continue
												} else {
													*v -= *v >> kNumMoveBits
													rCode -= bound
													rRange -= bound

													// Normalize
													if rRange < kTopValue {
														b, err := r.rangeDec.inStream.ReadByte()
														if err != nil {
															return err
														}

														rRange <<= 8
														rCode = (rCode << 8) | uint32(b)
													}
												}
											}
										}
									} else {
										*v -= *v >> kNumMoveBits
										rCode -= bound
										rRange -= bound

										// Normalize
										if rRange < kTopValue {
											b, err := r.rangeDec.inStream.ReadByte()
											if err != nil {
												return err
											}

											rRange <<= 8
											rCode = (rCode << 8) | uint32(b)
										}

										{ // rep match
											dist := uint32(0)

											{ // r.rangeDec.DecodeBit(&s.isRepG1[s.state])
												v := (*prob)(unsafe.Pointer(uintptr(unsafe.Pointer(&s.isRepG1[0])) + uintptr(s.state)*unsafe.Sizeof(prob(0))))
												bound := (rRange >> kNumBitModelTotalBits) * uint32(*v)

												if rCode < bound {
													*v += ((1 << kNumBitModelTotalBits) - *v) >> kNumMoveBits
													rRange = bound
													dist = s.rep1
													s.rep1 = s.rep0
													s.rep0 = dist

													// Normalize
													if rRange < kTopValue {
														b, err := r.rangeDec.inStream.ReadByte()
														if err != nil {
															return err
														}

														rRange <<= 8
														rCode = (rCode << 8) | uint32(b)
													}
												} else {
													*v -= *v >> kNumMoveBits
													rCode -= bound
													rRange -= bound

													// Normalize
													if rRange < kTopValue {
														b, err := r.rangeDec.inStream.ReadByte()
														if err != nil {
															return err
														}

														rRange <<= 8
														rCode = (rCode << 8) | uint32(b)
													}

													{ // isRepG1
														{ // r.rangeDec.DecodeBit(&s.isRepG2[s.state])
															v := (*prob)(unsafe.Pointer(uintptr(unsafe.Pointer(&s.isRepG2[0])) + uintptr(s.state)*unsafe.Sizeof(prob(0))))
															bound := (rRange >> kNumBitModelTotalBits) * uint32(*v)

															if rCode < bound {
																*v += ((1 << kNumBitModelTotalBits) - *v) >> kNumMoveBits
																rRange = bound

																dist = s.rep2
																s.rep2 = s.rep1
																s.rep1 = s.rep0
																s.rep0 = dist

																// Normalize
																if rRange < kTopValue {
																	b, err := r.rangeDec.inStream.ReadByte()
																	if err != nil {
																		return err
																	}

																	rRange <<= 8
																	rCode = (rCode << 8) | uint32(b)
																}
															} else {
																*v -= *v >> kNumMoveBits
																rCode -= bound
																rRange -= bound

																dist = s.rep3
																s.rep3 = s.rep2
																s.rep2 = s.rep1
																s.rep1 = s.rep0
																s.rep0 = dist

																// Normalize
																if rRange < kTopValue {
																	b, err := r.rangeDec.inStream.ReadByte()
																	if err != nil {
																		return err
																	}

																	rRange <<= 8
																	rCode = (rCode << 8) | uint32(b)
																}
															}
														}
													}
												}
											}
										}
									}
								}

								{ // r.s.repLenDecoder.Decode
									{ // r.rangeDec.DecodeBit(&r.s.repLenDecoder.choice)
										bound := (rRange >> kNumBitModelTotalBits) * uint32(s.repLenDecoderChoice)

										if rCode < bound {
											s.repLenDecoderChoice += ((1 << kNumBitModelTotalBits) - s.repLenDecoderChoice) >> kNumMoveBits
											rRange = bound

											// Normalize
											if rRange < kTopValue {
												b, err := r.rangeDec.inStream.ReadByte()
												if err != nil {
													return err
												}

												rRange <<= 8
												rCode = (rCode << 8) | uint32(b)
											}

											{ // r.s.repLenDecoder.lowCoder[s.posState].Decode
												m := uint32(1)

												probsPtr := uintptr(unsafe.Pointer(&s.repLenDecoderLowCoder[s.posState][0]))

												for i := 0; i < lenLowCoderNumBits; i++ {
													probPtr := (*prob)(unsafe.Pointer(probsPtr + uintptr(m)*unsafe.Sizeof(prob(0))))
													{ // rc.DecodeBit
														bound := (rRange >> kNumBitModelTotalBits) * uint32(*probPtr)

														if rCode < bound {
															*probPtr += ((1 << kNumBitModelTotalBits) - *probPtr) >> kNumMoveBits
															rRange = bound
															m <<= 1

															// Normalize
															if rRange < kTopValue {
																b, err := r.rangeDec.inStream.ReadByte()
																if err != nil {
																	return err
																}

																rRange <<= 8
																rCode = (rCode << 8) | uint32(b)
															}
														} else {
															*probPtr -= *probPtr >> kNumMoveBits
															rCode -= bound
															rRange -= bound
															m = (m << 1) | 1

															// Normalize
															if rRange < kTopValue {
																b, err := r.rangeDec.inStream.ReadByte()
																if err != nil {
																	return err
																}

																rRange <<= 8
																rCode = (rCode << 8) | uint32(b)
															}
														}
													}
												}

												s.state = stateUpdateRep(s.state)
												length = m - (uint32(1) << lenLowCoderNumBits) + kMatchMinLen

												if s.unpackSizeDefined && uint32(s.bytesLeft) < length {
													length = uint32(s.bytesLeft)
													r.outWindow.CopyMatch(s.rep0+1, length)
													s.bytesLeft -= uint64(length)

													return ErrResultError
												} else {
													r.outWindow.CopyMatch(s.rep0+1, length)
													s.bytesLeft -= uint64(length)

													continue
												}
											}
										} else {
											s.repLenDecoderChoice -= s.repLenDecoderChoice >> kNumMoveBits
											rCode -= bound
											rRange -= bound

											// Normalize
											if rRange < kTopValue {
												b, err := r.rangeDec.inStream.ReadByte()
												if err != nil {
													return err
												}

												rRange <<= 8
												rCode = (rCode << 8) | uint32(b)
											}

											{ // r.rangeDec.DecodeBit(&r.s.repLenDecoder.choice2)
												bound := (rRange >> kNumBitModelTotalBits) * uint32(s.repLenDecoderChoice2)

												if rCode < bound {
													s.repLenDecoderChoice2 += ((1 << kNumBitModelTotalBits) - s.repLenDecoderChoice2) >> kNumMoveBits
													rRange = bound

													// Normalize
													if rRange < kTopValue {
														b, err := r.rangeDec.inStream.ReadByte()
														if err != nil {
															return err
														}

														rRange <<= 8
														rCode = (rCode << 8) | uint32(b)
													}

													{ // r.s.repLenDecoder.midCoder[s.posState].Decode
														m := uint32(1)

														probsPtr := uintptr(unsafe.Pointer(&s.repLenDecoderMidCoder[s.posState][0]))

														for i := 0; i < lenMidCoderNumBits; i++ {
															probPtr := (*prob)(unsafe.Pointer(probsPtr + uintptr(m)*unsafe.Sizeof(prob(0))))
															{ // rc.DecodeBit
																bound := (rRange >> kNumBitModelTotalBits) * uint32(*probPtr)

																if rCode < bound {
																	*probPtr += ((1 << kNumBitModelTotalBits) - *probPtr) >> kNumMoveBits
																	rRange = bound
																	m <<= 1

																	// Normalize
																	if rRange < kTopValue {
																		b, err := r.rangeDec.inStream.ReadByte()
																		if err != nil {
																			return err
																		}

																		rRange <<= 8
																		rCode = (rCode << 8) | uint32(b)
																	}
																} else {
																	*probPtr -= *probPtr >> kNumMoveBits
																	rCode -= bound
																	rRange -= bound
																	m = (m << 1) | 1

																	// Normalize
																	if rRange < kTopValue {
																		b, err := r.rangeDec.inStream.ReadByte()
																		if err != nil {
																			return err
																		}

																		rRange <<= 8
																		rCode = (rCode << 8) | uint32(b)
																	}
																}
															}
														}

														s.state = stateUpdateRep(s.state)
														length = 8 + m - (uint32(1) << lenMidCoderNumBits) + kMatchMinLen

														if s.unpackSizeDefined && uint32(s.bytesLeft) < length {
															length = uint32(s.bytesLeft)
															r.outWindow.CopyMatch(s.rep0+1, length)
															s.bytesLeft -= uint64(length)

															return ErrResultError
														} else {
															r.outWindow.CopyMatch(s.rep0+1, length)
															s.bytesLeft -= uint64(length)

															continue
														}
													}
												} else {
													s.repLenDecoderChoice2 -= s.repLenDecoderChoice2 >> kNumMoveBits
													rCode -= bound
													rRange -= bound

													// Normalize
													if rRange < kTopValue {
														b, err := r.rangeDec.inStream.ReadByte()
														if err != nil {
															return err
														}

														rRange <<= 8
														rCode = (rCode << 8) | uint32(b)
													}

													{ // r.s.repLenDecoder.highCoder.Decode
														m := uint32(1)

														probsPtr := uintptr(unsafe.Pointer(&s.repLenDecoderHighCoder[0]))

														for i := 0; i < lenHighCoderNumBits; i++ {
															probPtr := (*prob)(unsafe.Pointer(probsPtr + uintptr(m)*unsafe.Sizeof(prob(0))))
															{ // rc.DecodeBit
																bound := (rRange >> kNumBitModelTotalBits) * uint32(*probPtr)

																if rCode < bound {
																	*probPtr += ((1 << kNumBitModelTotalBits) - *probPtr) >> kNumMoveBits
																	rRange = bound
																	m <<= 1

																	// Normalize
																	if rRange < kTopValue {
																		b, err := r.rangeDec.inStream.ReadByte()
																		if err != nil {
																			return err
																		}

																		rRange <<= 8
																		rCode = (rCode << 8) | uint32(b)
																	}
																} else {
																	*probPtr -= *probPtr >> kNumMoveBits
																	rCode -= bound
																	rRange -= bound
																	m = (m << 1) | 1

																	// Normalize
																	if rRange < kTopValue {
																		b, err := r.rangeDec.inStream.ReadByte()
																		if err != nil {
																			return err
																		}

																		rRange <<= 8
																		rCode = (rCode << 8) | uint32(b)
																	}
																}
															}
														}

														s.state = stateUpdateRep(s.state)
														length = 16 + m - (uint32(1) << lenHighCoderNumBits) + kMatchMinLen

														if s.unpackSizeDefined && uint32(s.bytesLeft) < length {
															length = uint32(s.bytesLeft)
															r.outWindow.CopyMatch(s.rep0+1, length)
															s.bytesLeft -= uint64(length)

															return ErrResultError
														} else {
															r.outWindow.CopyMatch(s.rep0+1, length)
															s.bytesLeft -= uint64(length)

															continue
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	r.rangeDec.Code = rCode
	r.rangeDec.Range = rRange

	return
}
