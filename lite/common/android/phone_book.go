package android

import "CentralizedControl/common/utils"

type PhoneBookItem struct {
	FirstName  string
	SecondName string
	Email      []string
	Phone      []string
}

func GenEmail(randTool *utils.RandTool) string {
	return randTool.GenString(utils.CharSet_abc, 8) + "@gmail.com"
}

func GenPhoneNumber(randTool *utils.RandTool, country string) (string, string) {
	sims := Resource.deviceResource.Sim[country]
	sim := utils.ChoseOne2(randTool, sims)
	perf := utils.ChoseOne2(randTool, sim.PhoneNumberPref)
	return sim.AreaCode, perf + randTool.GenString(utils.CharSet_123, sim.PhoneNumberLength-len(perf))
}

func createPhoneBookItem(randTool *utils.RandTool, country string) PhoneBookItem {
	pb := PhoneBookItem{}
	pc := randTool.GenNumber(1, 2)
	pb.Phone = make([]string, pc)
	for i := 0; i < pc; i++ {
		_, pb.Phone[i] = GenPhoneNumber(randTool, country)
	}
	ec := randTool.GenNumber(0, 1)
	pb.Email = make([]string, ec)
	for i := 0; i < ec; i++ {
		pb.Email[i] = GenEmail(randTool)
	}
	pb.FirstName = ChoiceFirstName(randTool)
	pb.SecondName = ChoiceSecondName(randTool)
	return pb
}

type PhoneBook struct {
	Item []PhoneBookItem
}

func CreatePhoneBook(randTool *utils.RandTool, country string) *PhoneBook {
	count := randTool.GenNumber(5, 20)
	pb := &PhoneBook{
		Item: make([]PhoneBookItem, count),
	}
	for i := 0; i < count; i++ {
		pb.Item[i] = createPhoneBookItem(randTool, country)
	}
	return pb
}
