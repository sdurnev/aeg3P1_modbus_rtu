package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/goburrow/modbus"
	"os"
	"time"
)

/*
!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
!!!!!!!!!!!! VERSION !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
*/
const version = "0.0.1"

/*
!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
!!!!!!!!!!!! VERSION !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
*/

type par struct {
	reg     uint16
	regtype uint
	regName string
}
type param []par

var aegPar = param{
	{1, 0, "01_UPSStatus"},
	{2, 1, "02_NonCriticalFault"},
	{3, 1, "03_CriticalFault"},
	{4, 1, "04_Userdef_Digln1Аctive"},
	{5, 1, "05_Userdef_Digln2Аctive"},
	{6, 1, "06_Userdef_Digln3Аctive"},
	{7, 1, "07_EmergencySwitchOff"},
	{8, 1, "08_DouCanFault"},
	{9, 1, "09_MainsFault"},
	{10, 1, "10_RectifierWarning"},
	{11, 1, "11_RectifierFault"},
	{12, 1, "12_BatteryАlarm"},
	{13, 1, "13_lnverterWarning"},
	{14, 1, "14_lnverterFault"},
	{15, 1, "15_SBSWarning"},
	{16, 1, "16_SBSFault"},
	{17, 1, "17_EqualisingCharge"},
	{18, 1, "18_Charge"},
	{19, 1, "19_TrickleCharge"},
	{20, 1, "20_GenSetOperation"},
	{21, 1, "21_BattTemp.SensFault"},
	{22, 1, "22_BatteryTemp.ToHigh"},
	{23, 1, "23_CircuitBreaker"},
	{24, 1, "24_BatteryWarning"},
	{25, 1, "25_Batterylow"},
	{26, 1, "26_Warning_InverterOverload"},
	{27, 1, "27_FanFault"},
	{28, 1, "28_Alarm_InverterOverload"},
	{29, 1, "29_ShortCircuit"},
	{30, 1, "30_DcUnderVoltage"},
	{31, 1, "31_DcOverVoltage"},
	{32, 1, "32_PowerStackOverТemp"},
	{33, 1, "33_SBSReady"},
	{34, 1, "34_SBSMainsFault"},
	{35, 1, "35_SBSBlocked"},
	{36, 1, "36_RectifierOn"},
	{37, 1, "37_1nverterOn"},
	{38, 1, "38_SBSOn"},
	{39, 2, "39_RectMainsFreq."},
	{40, 1, "40_RectMainsVoltL1"},
	{41, 1, "41_RectMainsVoltL2"},
	{42, 1, "42_RectMainsVoltLЗ"},
	{43, 2, "43_SBSMainsFreq."},
	{44, 1, "44_SBSMainsVoltL1"},
	{45, 1, "45_BatteryVoltage"},
	{46, 3, "46_BatteryCurrent"},
	{47, 2, "47_АutonomyTime"},
	{48, 1, "48_BatteryCapacity"},
	{49, 3, "49_BatteryTemperature"},
	{50, 2, "50_OutputFreq."},
	{51, 1, "51_OutputVoltageL1"},
	{52, 1, "52_OutputLoadL1"},
	{53, 1, "53_OutputCurrentL1"},
	{54, 2, "54_OutputPowerL1"},
	{55, 1, "55_LifeCheck"},
	{56, 1, "56_Userdef_AUX1-Rect."},
	{57, 1, "57_Userdef_AUX2-Rect."},
	{58, 1, "58_Userdef_AUXЗ-Rect."},
	{59, 1, "59_Userdef_AUX4-Rect."},
	{60, 1, "60_Userdef_AUX5-Rect."},
	{61, 1, "61_Userdef_AUX6-Rect."},
	{62, 1, "62_Userdef_AUX7-Rect."},
	{63, 1, "63_Userdef_AUX1-lnv."},
	{64, 1, "64_Userdef_AUX2-lnv."},
	{65, 1, "65_Userdef_AUXЗ-lnv."},
	{66, 1, "66_Userdef_AUX4-lnv."},
	{67, 1, "67_Userdef_AUX5-lnv."},
	{68, 1, "68_Userdef_AUX6-lnv."},
	{69, 1, "69_Userdef_AUX7-lnv."},
	{70, 1, "70_Userdef_AUX1-SBS"},
	{71, 1, "71_Userdef_AUX2-SBS"},
	{72, 1, "72_Userdef_AUXЗ-SBS"},
	{73, 1, "73_Userdef_AUX4-SBS"},
	{74, 1, "74_Userdef_AUX5-SBS"},
	{75, 1, "75_Userdef_AUX6-SB"},
	{76, 1, "76_Userdef_AUX7-SBS"},
}

type respStruct struct {
	namePar  string
	valuePar string
}

var respons []respStruct

func main() {

	serialPort := flag.String("serial", "/dev/ttyRS485-1", "a string")
	serialSpeed := flag.Int("speed", 9600, "a int")
	slaveID := flag.Int("id", 1, "an int")
	timeout := flag.Int("t", 3000, "an int")
	typOfReg := flag.Int("type", 1, "an int")
	regQuantity := flag.Uint("q", 76, "an uint")
	/*dataBits := flag.Int("dbits", 1, "a int")
	parity := flag.String("parity", "E", "a string")
	stopBits := flag.Int("sbits", 1, "an int")*/
	flag.Parse()

	handler := modbus.NewRTUClientHandler(fmt.Sprint(*serialPort))
	handler.BaudRate = *serialSpeed
	handler.SlaveId = byte(*slaveID)
	handler.Timeout = time.Duration(*timeout) * time.Millisecond
	/*handler.DataBits = int(*dataBits)
	handler.Parity = fmt.Sprint(*parity)
	handler.StopBits = *stopBits*/

	defer handler.Close()
	client := modbus.NewClient(handler)

	/**
	Конфигурация для чтения параметров из json файла ./config.json
	*/
	/*		data, err := ioutil.ReadFile("./config.json")
			if err != nil {
				fmt.Print(err)
			}
			type Config struct {
				DEV       string
				SPEED     int
				DATABITS  int
				PARITY    string
				STOPBITS  int
				TIMEOUT   time.Duration
				ID        int
				STARTREG  int
				QUANTITY  int
				TYPEOFREG int
			}
			var conf Config
			err = json.Unmarshal(data, &conf)
			if err != nil {
				printError(err)
			}

			handler := modbus.NewRTUClientHandler(conf.DEV)
			handler.BaudRate = conf.SPEED
			handler.DataBits = conf.DATABITS
			handler.Parity = conf.PARITY
			handler.StopBits = conf.STOPBITS
			handler.SlaveId = byte(conf.ID)
			handler.Timeout = conf.TIMEOUT * time.Millisecond
	*/
	if *typOfReg == 1 {
		results, err := client.ReadHoldingRegisters(uint16(1), uint16(*regQuantity))
		if err != nil {
			printError(err)
		}
		printResult(results)
		/*fmt.Println(len(results))
		fmt.Println(results)
		os.Exit(0)*/
	}

	if *typOfReg == 2 {
		results, err := client.ReadInputRegisters(uint16(1), uint16(*regQuantity))
		if err != nil {
			printError(err)
		}
		printResult(results)
		/*fmt.Println(len(results))
		fmt.Println(results)
		os.Exit(0)*/
	}
	fmt.Printf("{\"status\":\"error\", \"error\":\"typeofreg not 1 or 2 \"} \n")

	/*var ttt = []byte{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 1, 244, 0, 226, 0, 225, 0, 225, 1, 244, 0, 230, 0, 229, 0, 229, 0, 243, 0, 0, 5, 220, 0, 100, 4, 86, 1, 244, 0, 219, 0, 6, 0, 1, 0, 1, 0, 219, 0, 19, 0, 2, 0, 1, 0, 219, 0, 43, 0, 4, 0, 4, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	var ttt = []byte{}
	printResult(ttt)*/
}

func printError(err error) {
	fmt.Printf("{\"status\":\"error\", \"error\":\"%s\", \"version\":\"%s\"} \n", err, version)
	os.Exit(1)
}

func printResult(data []byte) {
	if len(data) != 0 {
		for i := 0; i < (len(data) / 2); i++ {
			if aegPar[i].regtype == 0 {
				temp1 := aegPar[i].regName
				var temp2 string
				a := binary.BigEndian.Uint16(data[i*2 : (i*2)+2])
				switch a {
				case 1:
					temp2 = "\"Normal_mode\""
				case 2:
					temp2 = "\"Baypass_mode\""
				case 3:
					temp2 = "\"Battery_mode\""
				case 4:
					temp2 = "\"Eco_mode\""
				case 6:
					temp2 = "\"Shutdown_run\""
				case 7:
					temp2 = "\"Shutdown\""
				}
				r := respStruct{namePar: temp1, valuePar: temp2}
				respons = append(respons, r)
			} else if aegPar[i].regtype == 2 {
				temp1 := aegPar[i].regName
				data1 := float64(binary.BigEndian.Uint16(data[i*2:(i*2)+2])) / 10
				temp2 := fmt.Sprintf("%.2f", data1)
				r := respStruct{namePar: temp1, valuePar: temp2}
				respons = append(respons, r)
			} else if aegPar[i].regtype == 3 {
				temp1 := aegPar[i].regName
				i := binary.BigEndian.Uint16(data[i*2 : (i*2)+2])
				b := float64(int16(i)) / 10
				temp2 := fmt.Sprintf("%.2f", b)
				r := respStruct{namePar: temp1, valuePar: temp2}
				respons = append(respons, r)
			} else {
				temp1 := aegPar[i].regName
				data1 := binary.BigEndian.Uint16(data[i*2 : (i*2)+2])
				temp2 := fmt.Sprintf("%d", data1)
				r := respStruct{namePar: temp1, valuePar: temp2}
				respons = append(respons, r)
			}
		}
		printJson(respons)
	} else {
		fmt.Printf("{\"status\":\"error\", \"error\":\"lengt of data is 0\", \"version\":\"%s\"} \n", version)
		os.Exit(100)
	}

}

func printJson(data []respStruct) {
	for l := 0; l < len(data); l++ {
		if l == 0 {
			fmt.Printf("{")
		}
		fmt.Printf("\"%s\":", data[l].namePar)
		fmt.Printf("%s,", data[l].valuePar)
		if l == len(data)-1 {
			fmt.Printf("\"version\":\"%s\"}\n", version)
		}
	}
	os.Exit(0)
}

/* build for rapberry
env GOOS=linux GOARCH=arm GOARM=5 go build
*/
