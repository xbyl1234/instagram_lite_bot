package lzma

func initProbs(probs []prob) {
	for i := 0; i < len(probs); i++ {
		probs[i] = probInitVal
	}
}
