package fspdriver

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	MZI_N_NODES int = 64
	MMI_N_NODES int = MZI_N_NODES * 3
)

var (
	// MZI to MMI as of physical layout
	// See dye documentation for details
	MZI_MMI_GRID_MAP [64][3][]int = [64][3][]int{
		{{13, 14}, {15, 14}, {17, 14}}, // P0 (1) cba p
		{{19, 14}, {21, 14}, {23, 14}}, // P1 (2) cba p
		{{12, 15}, {14, 15}, {16, 15}}, // P2 (3) cba i
		{{18, 15}, {20, 15}, {22, 15}}, // P3 (4) cba i

		{{17, 12}, {15, 12}, {13, 12}}, // O0 (5) abc p
		{{23, 12}, {21, 12}, {19, 12}}, // 01 (6) abc p
		{{16, 13}, {14, 13}, {12, 13}}, // 02 (7) abc i
		{{22, 13}, {20, 13}, {18, 13}}, // O3 (8) abc i

		{{13, 10}, {15, 10}, {17, 10}}, // N0 (9) cba p
		{{19, 10}, {21, 10}, {23, 10}}, // N1 (10) cba p
		{{12, 11}, {14, 11}, {16, 11}}, // N2 (11) cba i
		{{18, 11}, {20, 11}, {22, 11}}, // N3 (12) cba i

		{{17, 8}, {15, 8}, {13, 8}}, // M0 (13) abc p
		{{23, 8}, {21, 8}, {19, 8}}, // M1 (14) abc p
		{{16, 9}, {14, 9}, {12, 9}}, // M2 (15) abc i
		{{22, 9}, {20, 9}, {18, 9}}, // M3 (16) abc i

		{{13, 6}, {15, 6}, {17, 6}}, // L0 (17) cba p
		{{19, 6}, {21, 6}, {23, 6}}, // L1 (18) cba p
		{{12, 7}, {14, 7}, {16, 7}}, // L2 (19) cba i
		{{18, 7}, {20, 7}, {22, 7}}, // L3 (20) cba i

		{{17, 4}, {15, 4}, {13, 4}}, // K0 (21) abc p
		{{23, 4}, {21, 4}, {19, 4}}, // K1 (22) abc p
		{{16, 5}, {14, 5}, {12, 5}}, // K2 (23) abc i
		{{22, 5}, {20, 5}, {18, 5}}, // K3 (24) abc i

		{{13, 2}, {15, 2}, {17, 2}}, // J0 (25) cba p
		{{19, 2}, {21, 2}, {23, 2}}, // J1 (26) cba p
		{{12, 3}, {14, 3}, {16, 3}}, // J2 (27) cba i
		{{18, 3}, {20, 3}, {22, 3}}, // J3 (28) cba i

		{{17, 0}, {15, 0}, {13, 0}}, // I0 (29) abc p
		{{23, 0}, {21, 0}, {19, 0}}, // I1 (30) abc p
		{{16, 1}, {14, 1}, {12, 1}}, // I2 (31) abc i
		{{22, 1}, {20, 1}, {18, 1}}, // I3 (32) abc i

		//

		{{0, 1}, {2, 1}, {4, 1}},  // H0 (33) cba p
		{{6, 1}, {8, 1}, {10, 1}}, // H1 (34) cba p
		{{1, 0}, {3, 0}, {5, 0}},  // H2 (35) cba i
		{{7, 0}, {9, 0}, {11, 0}}, // H3 (36) cba i

		{{4, 3}, {2, 3}, {0, 3}},  // G0 (37) abc p
		{{10, 3}, {8, 3}, {6, 3}}, // G1 (38) abc p
		{{5, 2}, {3, 2}, {1, 2}},  // G2 (39) abc i
		{{11, 2}, {9, 2}, {7, 2}}, // G3 (40) abc i

		{{0, 5}, {2, 5}, {4, 5}},  // F0 (41) cba p
		{{6, 5}, {8, 5}, {10, 5}}, // F1 (42) cba p
		{{1, 4}, {3, 4}, {5, 4}},  // F2 (43) cba i
		{{7, 4}, {9, 4}, {11, 4}}, // F3 (44) cba i

		{{4, 7}, {2, 7}, {0, 7}},  // E0 (45) abc p
		{{10, 7}, {8, 7}, {6, 7}}, // E1 (46) abc p
		{{5, 6}, {3, 6}, {1, 6}},  // E2 (47) abc i
		{{11, 6}, {9, 6}, {7, 6}}, // E3 (48) abc i

		{{0, 9}, {2, 9}, {4, 9}},  // D0 (49) cba p
		{{6, 9}, {8, 9}, {10, 9}}, // D1 (50) cba p
		{{1, 8}, {3, 8}, {5, 8}},  // D2 (51) cba i
		{{7, 8}, {9, 8}, {11, 8}}, // D3 (52) cba i

		{{4, 11}, {2, 11}, {0, 11}},  // C0 (53) abc p
		{{10, 11}, {8, 11}, {6, 11}}, // C1 (54) abc p
		{{5, 10}, {3, 10}, {1, 10}},  // C2 (55) abc i
		{{11, 10}, {9, 10}, {7, 10}}, // C3 (56) abc p

		{{0, 13}, {2, 13}, {4, 13}},  // B0 (57) cba p
		{{6, 13}, {8, 13}, {10, 13}}, // B1 (58) cba p
		{{1, 12}, {3, 12}, {5, 12}},  // B2 (59) cba i
		{{7, 12}, {9, 12}, {11, 12}}, // B3 (60) cba i

		{{4, 15}, {2, 15}, {0, 15}},  // A0 (61) abc p
		{{10, 15}, {8, 15}, {6, 15}}, // A1 (62) abc p
		{{5, 14}, {3, 14}, {1, 14}},  // A2 (63) abc i
		{{11, 14}, {9, 14}, {7, 14}}, // A3 (64) abc i
	}

	// Indexing: col-major, with interlacing rows
	// 12 - number of interlaced rows
	// row is int-divided by 2 because
	// the MZI_MMI_MAP uses deinterlaced row indexing
	//
	// To convert grid to flat index:
	// aIdx := a[1]*12 + a[0]/2
	// bIdx := b[1]*12 + b[0]/2
	// cIdx := c[1]*12 + c[0]/2
	MZI_MMI_INDICES_MAP = [64][3]int{
		{174, 175, 176},
		{177, 178, 179},
		{186, 187, 188},
		{189, 190, 191},
		{152, 151, 150},
		{155, 154, 153},
		{164, 163, 162},
		{167, 166, 165},
		{126, 127, 128},
		{129, 130, 131},
		{138, 139, 140},
		{141, 142, 143},
		{104, 103, 102},
		{107, 106, 105},
		{116, 115, 114},
		{119, 118, 117},
		{78, 79, 80},
		{81, 82, 83},
		{90, 91, 92},
		{93, 94, 95},
		{56, 55, 54},
		{59, 58, 57},
		{68, 67, 66},
		{71, 70, 69},
		{30, 31, 32},
		{33, 34, 35},
		{42, 43, 44},
		{45, 46, 47},
		{8, 7, 6},
		{11, 10, 9},
		{20, 19, 18},
		{23, 22, 21},
		{12, 13, 14},
		{15, 16, 17},
		{0, 1, 2},
		{3, 4, 5},
		{38, 37, 36},
		{41, 40, 39},
		{26, 25, 24},
		{29, 28, 27},
		{60, 61, 62},
		{63, 64, 65},
		{48, 49, 50},
		{51, 52, 53},
		{86, 85, 84},
		{89, 88, 87},
		{74, 73, 72},
		{77, 76, 75},
		{108, 109, 110},
		{111, 112, 113},
		{96, 97, 98},
		{99, 100, 101},
		{134, 133, 132},
		{137, 136, 135},
		{122, 121, 120},
		{125, 124, 123},
		{156, 157, 158},
		{159, 160, 161},
		{144, 145, 146},
		{147, 148, 149},
		{182, 181, 180},
		{185, 184, 183},
		{170, 169, 168},
		{173, 172, 171},
	}
)


const (
	CAMERA_GET_STATE_MQTT_TOPIC_PATH    = "/camera/state/get"
	CAMERA_GET_STATE_CB_MQTT_TOPIC_PATH = "/camera/state/get/cb"

	CAMERA_SET_STATE_MQTT_TOPIC_PATH    = "/camera/state/set"
	CAMERA_SET_STATE_CB_MQTT_TOPIC_PATH = "/camera/state/set/cb"

	CAMERA_GET_FRAMERATE_MQTT_TOPIC_PATH    = "/camera/framerate/get"
	CAMERA_GET_FRAMERATE_CB_MQTT_TOPIC_PATH = "/camera/framerate/get/cb"

	CAMERA_SET_FRAMERATE_MQTT_TOPIC_PATH    = "/camera/framerate/set"
	CAMERA_SET_FRAMERATE_CB_MQTT_TOPIC_PATH = "/camera/framerate/set/cb"

	CAMERA_GET_CALIBRATION_MQTT_TOPIC_PATH    = "/camera/calibration/get"
	CAMERA_GET_CALIBRATION_CB_MQTT_TOPIC_PATH = "/camera/calibration/get/cb"

	// Calibration is performed each time CAMERA_STATE is turned from 0 to 1
	// No need to implement these callbacks
	// CAMERA_PERFORM_CALIBRATION_MQTT_TOPIC_PATH = "/camera/perform_calibration"
	// CAMERA_PERFORM_CALIBRATION_CB_MQTT_TOPIC_PATH = "/camera/perform_calibration/cb"

	CAMERA_GET_IMAGE_MQTT_TOPIC_PATH    = "/camera/get_image"
	CAMERA_GET_IMAGE_CB_MQTT_TOPIC_PATH = "/camera/get_image/cb"

	CAMERA_GET_DRAWING_MQTT_TOPIC_PATH    = "/camera/get_drawing"
	CAMERA_GET_DRAWING_CB_MQTT_TOPIC_PATH = "/camera/get_drawing/cb"

	CAMERA_MMI_BROADCAST_MQTT_TOPIC_PATH = "/camera/mmi/broadcast"
	CAMERA_MZI_BROADCAST_MQTT_TOPIC_PATH = "/camera/mzi/broadcast"
)

var (
	SEONE_SN_PATH = filepath.Join("config", "serialnumber.txt")
	SEONE_SN = ""
)

func InitSerialNumber() {
	sn, err := os.ReadFile(SEONE_SN_PATH)
	if err != nil {
		ERRORLogger.Fatal(err)
	}
	if len(sn) != 0 {
		if LOG_LEVEL <= INFO_LEVEL {
			INFOLogger.Printf("Setting SEONE_SN value: %s", string(sn))
		}
	}
	snStr := string(sn)
	snStr = strings.TrimSpace(snStr)
	SEONE_SN = snStr
	if SEONE_SN == "" {
		ERRORLogger.Fatal("Could not get seone's SN, exiting..")
	}
}
