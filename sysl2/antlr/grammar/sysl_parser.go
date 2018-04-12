// Generated from SyslParser.g4 by ANTLR 4.7.

package parser // SyslParser

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = reflect.Copy
var _ = strconv.Itoa

var parserATN = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 67, 665,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4, 18, 9,
	18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 4, 23, 9, 23,
	4, 24, 9, 24, 4, 25, 9, 25, 4, 26, 9, 26, 4, 27, 9, 27, 4, 28, 9, 28, 4,
	29, 9, 29, 4, 30, 9, 30, 4, 31, 9, 31, 4, 32, 9, 32, 4, 33, 9, 33, 4, 34,
	9, 34, 4, 35, 9, 35, 4, 36, 9, 36, 4, 37, 9, 37, 4, 38, 9, 38, 4, 39, 9,
	39, 4, 40, 9, 40, 4, 41, 9, 41, 4, 42, 9, 42, 4, 43, 9, 43, 4, 44, 9, 44,
	4, 45, 9, 45, 4, 46, 9, 46, 4, 47, 9, 47, 4, 48, 9, 48, 4, 49, 9, 49, 4,
	50, 9, 50, 4, 51, 9, 51, 4, 52, 9, 52, 4, 53, 9, 53, 4, 54, 9, 54, 4, 55,
	9, 55, 4, 56, 9, 56, 4, 57, 9, 57, 4, 58, 9, 58, 4, 59, 9, 59, 4, 60, 9,
	60, 4, 61, 9, 61, 4, 62, 9, 62, 4, 63, 9, 63, 4, 64, 9, 64, 4, 65, 9, 65,
	4, 66, 9, 66, 4, 67, 9, 67, 4, 68, 9, 68, 3, 2, 3, 2, 3, 2, 3, 2, 7, 2,
	141, 10, 2, 12, 2, 14, 2, 144, 11, 2, 3, 3, 3, 3, 3, 3, 3, 3, 5, 3, 150,
	10, 3, 3, 3, 3, 3, 3, 4, 3, 4, 3, 4, 7, 4, 157, 10, 4, 12, 4, 14, 4, 160,
	11, 4, 3, 5, 3, 5, 3, 5, 3, 5, 3, 6, 3, 6, 3, 6, 3, 6, 3, 7, 3, 7, 3, 7,
	3, 8, 3, 8, 3, 9, 3, 9, 3, 9, 3, 9, 7, 9, 179, 10, 9, 12, 9, 14, 9, 182,
	11, 9, 3, 9, 3, 9, 3, 10, 3, 10, 3, 10, 3, 10, 3, 11, 3, 11, 3, 11, 3,
	11, 3, 11, 5, 11, 195, 10, 11, 3, 12, 3, 12, 3, 12, 3, 12, 7, 12, 201,
	10, 12, 12, 12, 14, 12, 204, 11, 12, 3, 12, 3, 12, 3, 13, 3, 13, 5, 13,
	210, 10, 13, 3, 14, 3, 14, 3, 14, 3, 14, 7, 14, 216, 10, 14, 12, 14, 14,
	14, 219, 11, 14, 3, 14, 3, 14, 3, 15, 3, 15, 3, 15, 5, 15, 226, 10, 15,
	3, 16, 3, 16, 3, 17, 3, 17, 3, 18, 3, 18, 3, 18, 6, 18, 235, 10, 18, 13,
	18, 14, 18, 236, 3, 18, 3, 18, 3, 19, 3, 19, 3, 19, 5, 19, 244, 10, 19,
	3, 20, 3, 20, 3, 20, 3, 20, 3, 20, 3, 21, 3, 21, 6, 21, 253, 10, 21, 13,
	21, 14, 21, 254, 3, 21, 3, 21, 3, 22, 3, 22, 3, 22, 3, 22, 5, 22, 263,
	10, 22, 3, 22, 5, 22, 266, 10, 22, 3, 22, 5, 22, 269, 10, 22, 3, 22, 5,
	22, 272, 10, 22, 3, 22, 3, 22, 5, 22, 276, 10, 22, 5, 22, 278, 10, 22,
	3, 23, 3, 23, 3, 23, 3, 23, 3, 24, 7, 24, 285, 10, 24, 12, 24, 14, 24,
	288, 11, 24, 3, 24, 3, 24, 3, 24, 5, 24, 293, 10, 24, 3, 24, 3, 24, 3,
	24, 3, 24, 6, 24, 299, 10, 24, 13, 24, 14, 24, 300, 3, 24, 3, 24, 3, 25,
	3, 25, 3, 26, 3, 26, 3, 26, 3, 27, 3, 27, 7, 27, 312, 10, 27, 12, 27, 14,
	27, 315, 11, 27, 3, 28, 3, 28, 5, 28, 319, 10, 28, 3, 28, 5, 28, 322, 10,
	28, 3, 29, 3, 29, 3, 29, 3, 30, 3, 30, 3, 30, 3, 30, 5, 30, 331, 10, 30,
	6, 30, 333, 10, 30, 13, 30, 14, 30, 334, 3, 30, 3, 30, 3, 31, 3, 31, 3,
	31, 5, 31, 342, 10, 31, 3, 32, 7, 32, 345, 10, 32, 12, 32, 14, 32, 348,
	11, 32, 3, 32, 3, 32, 3, 32, 3, 32, 6, 32, 354, 10, 32, 13, 32, 14, 32,
	355, 3, 32, 3, 32, 3, 33, 3, 33, 3, 33, 3, 33, 3, 33, 3, 33, 3, 34, 3,
	34, 3, 34, 3, 34, 5, 34, 370, 10, 34, 3, 34, 3, 34, 3, 35, 3, 35, 6, 35,
	376, 10, 35, 13, 35, 14, 35, 377, 3, 36, 3, 36, 3, 36, 3, 36, 3, 37, 3,
	37, 3, 37, 3, 37, 7, 37, 388, 10, 37, 12, 37, 14, 37, 391, 11, 37, 3, 38,
	3, 38, 6, 38, 395, 10, 38, 13, 38, 14, 38, 396, 3, 38, 5, 38, 400, 10,
	38, 3, 39, 3, 39, 5, 39, 404, 10, 39, 3, 39, 6, 39, 407, 10, 39, 13, 39,
	14, 39, 408, 3, 39, 5, 39, 412, 10, 39, 3, 39, 5, 39, 415, 10, 39, 3, 39,
	5, 39, 418, 10, 39, 3, 39, 3, 39, 3, 40, 3, 40, 3, 40, 3, 41, 3, 41, 5,
	41, 427, 10, 41, 3, 42, 3, 42, 3, 43, 3, 43, 3, 43, 5, 43, 434, 10, 43,
	3, 44, 3, 44, 3, 44, 3, 44, 3, 44, 7, 44, 441, 10, 44, 12, 44, 14, 44,
	444, 11, 44, 3, 44, 3, 44, 3, 45, 3, 45, 3, 45, 3, 45, 3, 45, 7, 45, 453,
	10, 45, 12, 45, 14, 45, 456, 11, 45, 3, 45, 5, 45, 459, 10, 45, 3, 46,
	3, 46, 3, 46, 3, 47, 3, 47, 3, 47, 3, 47, 7, 47, 468, 10, 47, 12, 47, 14,
	47, 471, 11, 47, 3, 47, 3, 47, 3, 48, 3, 48, 3, 49, 3, 49, 3, 49, 3, 49,
	3, 49, 6, 49, 482, 10, 49, 13, 49, 14, 49, 483, 3, 49, 3, 49, 3, 50, 7,
	50, 489, 10, 50, 12, 50, 14, 50, 492, 11, 50, 3, 51, 3, 51, 3, 51, 3, 51,
	6, 51, 498, 10, 51, 13, 51, 14, 51, 499, 3, 51, 3, 51, 3, 52, 3, 52, 3,
	52, 3, 52, 6, 52, 508, 10, 52, 13, 52, 14, 52, 509, 3, 52, 3, 52, 3, 53,
	3, 53, 3, 54, 3, 54, 3, 54, 3, 54, 3, 54, 3, 54, 3, 54, 3, 54, 3, 54, 3,
	54, 3, 54, 3, 54, 5, 54, 528, 10, 54, 3, 55, 3, 55, 5, 55, 532, 10, 55,
	3, 55, 5, 55, 535, 10, 55, 3, 55, 3, 55, 3, 55, 7, 55, 540, 10, 55, 12,
	55, 14, 55, 543, 11, 55, 3, 55, 3, 55, 3, 56, 3, 56, 3, 56, 5, 56, 550,
	10, 56, 3, 57, 3, 57, 3, 58, 3, 58, 3, 58, 3, 58, 3, 58, 6, 58, 559, 10,
	58, 13, 58, 14, 58, 560, 3, 58, 3, 58, 5, 58, 565, 10, 58, 5, 58, 567,
	10, 58, 3, 59, 3, 59, 3, 59, 5, 59, 572, 10, 59, 3, 60, 3, 60, 3, 60, 3,
	61, 3, 61, 3, 61, 3, 61, 3, 61, 6, 61, 582, 10, 61, 13, 61, 14, 61, 583,
	3, 61, 3, 61, 5, 61, 588, 10, 61, 3, 62, 3, 62, 3, 62, 5, 62, 593, 10,
	62, 3, 62, 3, 62, 3, 62, 3, 62, 6, 62, 599, 10, 62, 13, 62, 14, 62, 600,
	3, 62, 3, 62, 5, 62, 605, 10, 62, 3, 63, 3, 63, 3, 63, 3, 63, 3, 63, 3,
	63, 3, 63, 3, 63, 6, 63, 615, 10, 63, 13, 63, 14, 63, 616, 3, 63, 3, 63,
	3, 64, 7, 64, 622, 10, 64, 12, 64, 14, 64, 625, 11, 64, 3, 64, 3, 64, 3,
	64, 3, 64, 3, 65, 5, 65, 632, 10, 65, 3, 65, 3, 65, 3, 65, 7, 65, 637,
	10, 65, 12, 65, 14, 65, 640, 11, 65, 3, 66, 3, 66, 3, 66, 7, 66, 645, 10,
	66, 12, 66, 14, 66, 648, 11, 66, 3, 67, 6, 67, 651, 10, 67, 13, 67, 14,
	67, 652, 3, 68, 5, 68, 656, 10, 68, 3, 68, 6, 68, 659, 10, 68, 13, 68,
	14, 68, 660, 3, 68, 3, 68, 3, 68, 2, 2, 69, 2, 4, 6, 8, 10, 12, 14, 16,
	18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52,
	54, 56, 58, 60, 62, 64, 66, 68, 70, 72, 74, 76, 78, 80, 82, 84, 86, 88,
	90, 92, 94, 96, 98, 100, 102, 104, 106, 108, 110, 112, 114, 116, 118, 120,
	122, 124, 126, 128, 130, 132, 134, 2, 7, 4, 2, 5, 5, 59, 59, 3, 2, 8, 9,
	4, 2, 58, 59, 63, 63, 4, 2, 31, 31, 59, 59, 4, 2, 59, 59, 63, 63, 2, 691,
	2, 136, 3, 2, 2, 2, 4, 145, 3, 2, 2, 2, 6, 153, 3, 2, 2, 2, 8, 161, 3,
	2, 2, 2, 10, 165, 3, 2, 2, 2, 12, 169, 3, 2, 2, 2, 14, 172, 3, 2, 2, 2,
	16, 174, 3, 2, 2, 2, 18, 185, 3, 2, 2, 2, 20, 189, 3, 2, 2, 2, 22, 196,
	3, 2, 2, 2, 24, 209, 3, 2, 2, 2, 26, 211, 3, 2, 2, 2, 28, 222, 3, 2, 2,
	2, 30, 227, 3, 2, 2, 2, 32, 229, 3, 2, 2, 2, 34, 231, 3, 2, 2, 2, 36, 243,
	3, 2, 2, 2, 38, 245, 3, 2, 2, 2, 40, 250, 3, 2, 2, 2, 42, 277, 3, 2, 2,
	2, 44, 279, 3, 2, 2, 2, 46, 286, 3, 2, 2, 2, 48, 304, 3, 2, 2, 2, 50, 306,
	3, 2, 2, 2, 52, 309, 3, 2, 2, 2, 54, 316, 3, 2, 2, 2, 56, 323, 3, 2, 2,
	2, 58, 326, 3, 2, 2, 2, 60, 338, 3, 2, 2, 2, 62, 346, 3, 2, 2, 2, 64, 359,
	3, 2, 2, 2, 66, 365, 3, 2, 2, 2, 68, 373, 3, 2, 2, 2, 70, 379, 3, 2, 2,
	2, 72, 383, 3, 2, 2, 2, 74, 394, 3, 2, 2, 2, 76, 411, 3, 2, 2, 2, 78, 421,
	3, 2, 2, 2, 80, 426, 3, 2, 2, 2, 82, 428, 3, 2, 2, 2, 84, 430, 3, 2, 2,
	2, 86, 435, 3, 2, 2, 2, 88, 447, 3, 2, 2, 2, 90, 460, 3, 2, 2, 2, 92, 463,
	3, 2, 2, 2, 94, 474, 3, 2, 2, 2, 96, 476, 3, 2, 2, 2, 98, 490, 3, 2, 2,
	2, 100, 493, 3, 2, 2, 2, 102, 503, 3, 2, 2, 2, 104, 513, 3, 2, 2, 2, 106,
	527, 3, 2, 2, 2, 108, 529, 3, 2, 2, 2, 110, 549, 3, 2, 2, 2, 112, 551,
	3, 2, 2, 2, 114, 566, 3, 2, 2, 2, 116, 571, 3, 2, 2, 2, 118, 573, 3, 2,
	2, 2, 120, 576, 3, 2, 2, 2, 122, 589, 3, 2, 2, 2, 124, 606, 3, 2, 2, 2,
	126, 623, 3, 2, 2, 2, 128, 631, 3, 2, 2, 2, 130, 641, 3, 2, 2, 2, 132,
	650, 3, 2, 2, 2, 134, 655, 3, 2, 2, 2, 136, 137, 7, 26, 2, 2, 137, 142,
	7, 59, 2, 2, 138, 139, 7, 25, 2, 2, 139, 141, 7, 59, 2, 2, 140, 138, 3,
	2, 2, 2, 141, 144, 3, 2, 2, 2, 142, 140, 3, 2, 2, 2, 142, 143, 3, 2, 2,
	2, 143, 3, 3, 2, 2, 2, 144, 142, 3, 2, 2, 2, 145, 146, 7, 44, 2, 2, 146,
	149, 7, 54, 2, 2, 147, 148, 7, 35, 2, 2, 148, 150, 7, 54, 2, 2, 149, 147,
	3, 2, 2, 2, 149, 150, 3, 2, 2, 2, 150, 151, 3, 2, 2, 2, 151, 152, 7, 45,
	2, 2, 152, 5, 3, 2, 2, 2, 153, 158, 5, 2, 2, 2, 154, 155, 7, 27, 2, 2,
	155, 157, 5, 2, 2, 2, 156, 154, 3, 2, 2, 2, 157, 160, 3, 2, 2, 2, 158,
	156, 3, 2, 2, 2, 158, 159, 3, 2, 2, 2, 159, 7, 3, 2, 2, 2, 160, 158, 3,
	2, 2, 2, 161, 162, 7, 40, 2, 2, 162, 163, 5, 6, 4, 2, 163, 164, 7, 41,
	2, 2, 164, 9, 3, 2, 2, 2, 165, 166, 7, 59, 2, 2, 166, 167, 7, 35, 2, 2,
	167, 168, 7, 59, 2, 2, 168, 11, 3, 2, 2, 2, 169, 170, 7, 49, 2, 2, 170,
	171, 7, 61, 2, 2, 171, 13, 3, 2, 2, 2, 172, 173, 7, 55, 2, 2, 173, 15,
	3, 2, 2, 2, 174, 175, 7, 40, 2, 2, 175, 180, 5, 14, 8, 2, 176, 177, 7,
	27, 2, 2, 177, 179, 5, 14, 8, 2, 178, 176, 3, 2, 2, 2, 179, 182, 3, 2,
	2, 2, 180, 178, 3, 2, 2, 2, 180, 181, 3, 2, 2, 2, 181, 183, 3, 2, 2, 2,
	182, 180, 3, 2, 2, 2, 183, 184, 7, 41, 2, 2, 184, 17, 3, 2, 2, 2, 185,
	186, 7, 40, 2, 2, 186, 187, 5, 16, 9, 2, 187, 188, 7, 41, 2, 2, 188, 19,
	3, 2, 2, 2, 189, 190, 7, 59, 2, 2, 190, 194, 7, 28, 2, 2, 191, 195, 5,
	14, 8, 2, 192, 195, 5, 16, 9, 2, 193, 195, 5, 18, 10, 2, 194, 191, 3, 2,
	2, 2, 194, 192, 3, 2, 2, 2, 194, 193, 3, 2, 2, 2, 195, 21, 3, 2, 2, 2,
	196, 197, 7, 40, 2, 2, 197, 202, 5, 20, 11, 2, 198, 199, 7, 27, 2, 2, 199,
	201, 5, 20, 11, 2, 200, 198, 3, 2, 2, 2, 201, 204, 3, 2, 2, 2, 202, 200,
	3, 2, 2, 2, 202, 203, 3, 2, 2, 2, 203, 205, 3, 2, 2, 2, 204, 202, 3, 2,
	2, 2, 205, 206, 7, 41, 2, 2, 206, 23, 3, 2, 2, 2, 207, 210, 5, 20, 11,
	2, 208, 210, 5, 2, 2, 2, 209, 207, 3, 2, 2, 2, 209, 208, 3, 2, 2, 2, 210,
	25, 3, 2, 2, 2, 211, 212, 7, 40, 2, 2, 212, 217, 5, 24, 13, 2, 213, 214,
	7, 27, 2, 2, 214, 216, 5, 24, 13, 2, 215, 213, 3, 2, 2, 2, 216, 219, 3,
	2, 2, 2, 217, 215, 3, 2, 2, 2, 217, 218, 3, 2, 2, 2, 218, 220, 3, 2, 2,
	2, 219, 217, 3, 2, 2, 2, 220, 221, 7, 41, 2, 2, 221, 27, 3, 2, 2, 2, 222,
	223, 7, 18, 2, 2, 223, 225, 9, 2, 2, 2, 224, 226, 5, 4, 3, 2, 225, 224,
	3, 2, 2, 2, 225, 226, 3, 2, 2, 2, 226, 29, 3, 2, 2, 2, 227, 228, 5, 28,
	15, 2, 228, 31, 3, 2, 2, 2, 229, 230, 7, 59, 2, 2, 230, 33, 3, 2, 2, 2,
	231, 232, 7, 33, 2, 2, 232, 234, 7, 3, 2, 2, 233, 235, 5, 12, 7, 2, 234,
	233, 3, 2, 2, 2, 235, 236, 3, 2, 2, 2, 236, 234, 3, 2, 2, 2, 236, 237,
	3, 2, 2, 2, 237, 238, 3, 2, 2, 2, 238, 239, 7, 4, 2, 2, 239, 35, 3, 2,
	2, 2, 240, 244, 7, 55, 2, 2, 241, 244, 5, 16, 9, 2, 242, 244, 5, 34, 18,
	2, 243, 240, 3, 2, 2, 2, 243, 241, 3, 2, 2, 2, 243, 242, 3, 2, 2, 2, 244,
	37, 3, 2, 2, 2, 245, 246, 7, 38, 2, 2, 246, 247, 7, 65, 2, 2, 247, 248,
	7, 28, 2, 2, 248, 249, 5, 36, 19, 2, 249, 39, 3, 2, 2, 2, 250, 252, 7,
	3, 2, 2, 251, 253, 5, 38, 20, 2, 252, 251, 3, 2, 2, 2, 253, 254, 3, 2,
	2, 2, 254, 252, 3, 2, 2, 2, 254, 255, 3, 2, 2, 2, 255, 256, 3, 2, 2, 2,
	256, 257, 7, 4, 2, 2, 257, 41, 3, 2, 2, 2, 258, 278, 5, 30, 16, 2, 259,
	263, 5, 10, 6, 2, 260, 263, 7, 5, 2, 2, 261, 263, 5, 32, 17, 2, 262, 259,
	3, 2, 2, 2, 262, 260, 3, 2, 2, 2, 262, 261, 3, 2, 2, 2, 263, 265, 3, 2,
	2, 2, 264, 266, 5, 4, 3, 2, 265, 264, 3, 2, 2, 2, 265, 266, 3, 2, 2, 2,
	266, 268, 3, 2, 2, 2, 267, 269, 7, 37, 2, 2, 268, 267, 3, 2, 2, 2, 268,
	269, 3, 2, 2, 2, 269, 271, 3, 2, 2, 2, 270, 272, 5, 26, 14, 2, 271, 270,
	3, 2, 2, 2, 271, 272, 3, 2, 2, 2, 272, 275, 3, 2, 2, 2, 273, 274, 7, 33,
	2, 2, 274, 276, 5, 40, 21, 2, 275, 273, 3, 2, 2, 2, 275, 276, 3, 2, 2,
	2, 276, 278, 3, 2, 2, 2, 277, 258, 3, 2, 2, 2, 277, 262, 3, 2, 2, 2, 278,
	43, 3, 2, 2, 2, 279, 280, 7, 59, 2, 2, 280, 281, 7, 22, 2, 2, 281, 282,
	5, 42, 22, 2, 282, 45, 3, 2, 2, 2, 283, 285, 7, 57, 2, 2, 284, 283, 3,
	2, 2, 2, 285, 288, 3, 2, 2, 2, 286, 284, 3, 2, 2, 2, 286, 287, 3, 2, 2,
	2, 287, 289, 3, 2, 2, 2, 288, 286, 3, 2, 2, 2, 289, 290, 9, 3, 2, 2, 290,
	292, 7, 59, 2, 2, 291, 293, 5, 26, 14, 2, 292, 291, 3, 2, 2, 2, 292, 293,
	3, 2, 2, 2, 293, 294, 3, 2, 2, 2, 294, 295, 7, 33, 2, 2, 295, 298, 7, 3,
	2, 2, 296, 299, 7, 57, 2, 2, 297, 299, 5, 44, 23, 2, 298, 296, 3, 2, 2,
	2, 298, 297, 3, 2, 2, 2, 299, 300, 3, 2, 2, 2, 300, 298, 3, 2, 2, 2, 300,
	301, 3, 2, 2, 2, 301, 302, 3, 2, 2, 2, 302, 303, 7, 4, 2, 2, 303, 47, 3,
	2, 2, 2, 304, 305, 9, 4, 2, 2, 305, 49, 3, 2, 2, 2, 306, 307, 7, 21, 2,
	2, 307, 308, 5, 48, 25, 2, 308, 51, 3, 2, 2, 2, 309, 313, 5, 48, 25, 2,
	310, 312, 5, 50, 26, 2, 311, 310, 3, 2, 2, 2, 312, 315, 3, 2, 2, 2, 313,
	311, 3, 2, 2, 2, 313, 314, 3, 2, 2, 2, 314, 53, 3, 2, 2, 2, 315, 313, 3,
	2, 2, 2, 316, 318, 5, 52, 27, 2, 317, 319, 7, 55, 2, 2, 318, 317, 3, 2,
	2, 2, 318, 319, 3, 2, 2, 2, 319, 321, 3, 2, 2, 2, 320, 322, 5, 26, 14,
	2, 321, 320, 3, 2, 2, 2, 321, 322, 3, 2, 2, 2, 322, 55, 3, 2, 2, 2, 323,
	324, 7, 59, 2, 2, 324, 325, 7, 33, 2, 2, 325, 57, 3, 2, 2, 2, 326, 327,
	7, 33, 2, 2, 327, 332, 7, 3, 2, 2, 328, 330, 7, 59, 2, 2, 329, 331, 5,
	26, 14, 2, 330, 329, 3, 2, 2, 2, 330, 331, 3, 2, 2, 2, 331, 333, 3, 2,
	2, 2, 332, 328, 3, 2, 2, 2, 333, 334, 3, 2, 2, 2, 334, 332, 3, 2, 2, 2,
	334, 335, 3, 2, 2, 2, 335, 336, 3, 2, 2, 2, 336, 337, 7, 4, 2, 2, 337,
	59, 3, 2, 2, 2, 338, 339, 9, 3, 2, 2, 339, 341, 7, 59, 2, 2, 340, 342,
	5, 58, 30, 2, 341, 340, 3, 2, 2, 2, 341, 342, 3, 2, 2, 2, 342, 61, 3, 2,
	2, 2, 343, 345, 7, 57, 2, 2, 344, 343, 3, 2, 2, 2, 345, 348, 3, 2, 2, 2,
	346, 344, 3, 2, 2, 2, 346, 347, 3, 2, 2, 2, 347, 349, 3, 2, 2, 2, 348,
	346, 3, 2, 2, 2, 349, 350, 7, 7, 2, 2, 350, 351, 5, 56, 29, 2, 351, 353,
	7, 3, 2, 2, 352, 354, 5, 60, 31, 2, 353, 352, 3, 2, 2, 2, 354, 355, 3,
	2, 2, 2, 355, 353, 3, 2, 2, 2, 355, 356, 3, 2, 2, 2, 356, 357, 3, 2, 2,
	2, 357, 358, 7, 4, 2, 2, 358, 63, 3, 2, 2, 2, 359, 360, 7, 38, 2, 2, 360,
	361, 7, 59, 2, 2, 361, 362, 7, 28, 2, 2, 362, 363, 7, 55, 2, 2, 363, 364,
	7, 56, 2, 2, 364, 65, 3, 2, 2, 2, 365, 366, 7, 30, 2, 2, 366, 369, 7, 42,
	2, 2, 367, 370, 5, 44, 23, 2, 368, 370, 7, 59, 2, 2, 369, 367, 3, 2, 2,
	2, 369, 368, 3, 2, 2, 2, 370, 371, 3, 2, 2, 2, 371, 372, 7, 43, 2, 2, 372,
	67, 3, 2, 2, 2, 373, 375, 7, 30, 2, 2, 374, 376, 9, 5, 2, 2, 375, 374,
	3, 2, 2, 2, 376, 377, 3, 2, 2, 2, 377, 375, 3, 2, 2, 2, 377, 378, 3, 2,
	2, 2, 378, 69, 3, 2, 2, 2, 379, 380, 7, 59, 2, 2, 380, 381, 7, 28, 2, 2,
	381, 382, 9, 2, 2, 2, 382, 71, 3, 2, 2, 2, 383, 384, 7, 37, 2, 2, 384,
	389, 5, 70, 36, 2, 385, 386, 7, 39, 2, 2, 386, 388, 5, 70, 36, 2, 387,
	385, 3, 2, 2, 2, 388, 391, 3, 2, 2, 2, 389, 387, 3, 2, 2, 2, 389, 390,
	3, 2, 2, 2, 390, 73, 3, 2, 2, 2, 391, 389, 3, 2, 2, 2, 392, 395, 5, 68,
	35, 2, 393, 395, 5, 66, 34, 2, 394, 392, 3, 2, 2, 2, 394, 393, 3, 2, 2,
	2, 395, 396, 3, 2, 2, 2, 396, 394, 3, 2, 2, 2, 396, 397, 3, 2, 2, 2, 397,
	399, 3, 2, 2, 2, 398, 400, 5, 72, 37, 2, 399, 398, 3, 2, 2, 2, 399, 400,
	3, 2, 2, 2, 400, 75, 3, 2, 2, 2, 401, 412, 5, 74, 38, 2, 402, 404, 7, 35,
	2, 2, 403, 402, 3, 2, 2, 2, 403, 404, 3, 2, 2, 2, 404, 406, 3, 2, 2, 2,
	405, 407, 7, 59, 2, 2, 406, 405, 3, 2, 2, 2, 407, 408, 3, 2, 2, 2, 408,
	406, 3, 2, 2, 2, 408, 409, 3, 2, 2, 2, 409, 412, 3, 2, 2, 2, 410, 412,
	7, 58, 2, 2, 411, 401, 3, 2, 2, 2, 411, 403, 3, 2, 2, 2, 411, 410, 3, 2,
	2, 2, 412, 414, 3, 2, 2, 2, 413, 415, 7, 55, 2, 2, 414, 413, 3, 2, 2, 2,
	414, 415, 3, 2, 2, 2, 415, 417, 3, 2, 2, 2, 416, 418, 5, 26, 14, 2, 417,
	416, 3, 2, 2, 2, 417, 418, 3, 2, 2, 2, 418, 419, 3, 2, 2, 2, 419, 420,
	7, 33, 2, 2, 420, 77, 3, 2, 2, 2, 421, 422, 7, 11, 2, 2, 422, 423, 7, 61,
	2, 2, 423, 79, 3, 2, 2, 2, 424, 427, 7, 35, 2, 2, 425, 427, 5, 52, 27,
	2, 426, 424, 3, 2, 2, 2, 426, 425, 3, 2, 2, 2, 427, 81, 3, 2, 2, 2, 428,
	429, 9, 6, 2, 2, 429, 83, 3, 2, 2, 2, 430, 433, 5, 80, 41, 2, 431, 432,
	7, 23, 2, 2, 432, 434, 5, 82, 42, 2, 433, 431, 3, 2, 2, 2, 433, 434, 3,
	2, 2, 2, 434, 85, 3, 2, 2, 2, 435, 436, 7, 12, 2, 2, 436, 437, 7, 63, 2,
	2, 437, 438, 7, 33, 2, 2, 438, 442, 7, 3, 2, 2, 439, 441, 5, 106, 54, 2,
	440, 439, 3, 2, 2, 2, 441, 444, 3, 2, 2, 2, 442, 440, 3, 2, 2, 2, 442,
	443, 3, 2, 2, 2, 443, 445, 3, 2, 2, 2, 444, 442, 3, 2, 2, 2, 445, 446,
	7, 4, 2, 2, 446, 87, 3, 2, 2, 2, 447, 458, 5, 86, 44, 2, 448, 449, 7, 13,
	2, 2, 449, 450, 7, 33, 2, 2, 450, 454, 7, 3, 2, 2, 451, 453, 5, 106, 54,
	2, 452, 451, 3, 2, 2, 2, 453, 456, 3, 2, 2, 2, 454, 452, 3, 2, 2, 2, 454,
	455, 3, 2, 2, 2, 455, 457, 3, 2, 2, 2, 456, 454, 3, 2, 2, 2, 457, 459,
	7, 4, 2, 2, 458, 448, 3, 2, 2, 2, 458, 459, 3, 2, 2, 2, 459, 89, 3, 2,
	2, 2, 460, 461, 7, 63, 2, 2, 461, 462, 7, 33, 2, 2, 462, 91, 3, 2, 2, 2,
	463, 464, 7, 14, 2, 2, 464, 465, 5, 90, 46, 2, 465, 469, 7, 3, 2, 2, 466,
	468, 5, 106, 54, 2, 467, 466, 3, 2, 2, 2, 468, 471, 3, 2, 2, 2, 469, 467,
	3, 2, 2, 2, 469, 470, 3, 2, 2, 2, 470, 472, 3, 2, 2, 2, 471, 469, 3, 2,
	2, 2, 472, 473, 7, 4, 2, 2, 473, 93, 3, 2, 2, 2, 474, 475, 7, 57, 2, 2,
	475, 95, 3, 2, 2, 2, 476, 477, 7, 16, 2, 2, 477, 478, 7, 63, 2, 2, 478,
	479, 7, 33, 2, 2, 479, 481, 7, 3, 2, 2, 480, 482, 5, 106, 54, 2, 481, 480,
	3, 2, 2, 2, 482, 483, 3, 2, 2, 2, 483, 481, 3, 2, 2, 2, 483, 484, 3, 2,
	2, 2, 484, 485, 3, 2, 2, 2, 485, 486, 7, 4, 2, 2, 486, 97, 3, 2, 2, 2,
	487, 489, 7, 59, 2, 2, 488, 487, 3, 2, 2, 2, 489, 492, 3, 2, 2, 2, 490,
	488, 3, 2, 2, 2, 490, 491, 3, 2, 2, 2, 491, 99, 3, 2, 2, 2, 492, 490, 3,
	2, 2, 2, 493, 494, 5, 98, 50, 2, 494, 495, 7, 33, 2, 2, 495, 497, 7, 3,
	2, 2, 496, 498, 5, 106, 54, 2, 497, 496, 3, 2, 2, 2, 498, 499, 3, 2, 2,
	2, 499, 497, 3, 2, 2, 2, 499, 500, 3, 2, 2, 2, 500, 501, 3, 2, 2, 2, 501,
	502, 7, 4, 2, 2, 502, 101, 3, 2, 2, 2, 503, 504, 7, 19, 2, 2, 504, 505,
	7, 33, 2, 2, 505, 507, 7, 3, 2, 2, 506, 508, 5, 100, 51, 2, 507, 506, 3,
	2, 2, 2, 508, 509, 3, 2, 2, 2, 509, 507, 3, 2, 2, 2, 509, 510, 3, 2, 2,
	2, 510, 511, 3, 2, 2, 2, 511, 512, 7, 4, 2, 2, 512, 103, 3, 2, 2, 2, 513,
	514, 7, 58, 2, 2, 514, 105, 3, 2, 2, 2, 515, 528, 5, 12, 7, 2, 516, 528,
	5, 88, 45, 2, 517, 528, 5, 92, 47, 2, 518, 528, 5, 78, 40, 2, 519, 528,
	5, 84, 43, 2, 520, 528, 5, 102, 52, 2, 521, 528, 5, 94, 48, 2, 522, 528,
	5, 96, 49, 2, 523, 528, 7, 55, 2, 2, 524, 528, 7, 17, 2, 2, 525, 528, 5,
	104, 53, 2, 526, 528, 5, 38, 20, 2, 527, 515, 3, 2, 2, 2, 527, 516, 3,
	2, 2, 2, 527, 517, 3, 2, 2, 2, 527, 518, 3, 2, 2, 2, 527, 519, 3, 2, 2,
	2, 527, 520, 3, 2, 2, 2, 527, 521, 3, 2, 2, 2, 527, 522, 3, 2, 2, 2, 527,
	523, 3, 2, 2, 2, 527, 524, 3, 2, 2, 2, 527, 525, 3, 2, 2, 2, 527, 526,
	3, 2, 2, 2, 528, 107, 3, 2, 2, 2, 529, 531, 7, 6, 2, 2, 530, 532, 5, 72,
	37, 2, 531, 530, 3, 2, 2, 2, 531, 532, 3, 2, 2, 2, 532, 534, 3, 2, 2, 2,
	533, 535, 5, 22, 12, 2, 534, 533, 3, 2, 2, 2, 534, 535, 3, 2, 2, 2, 535,
	536, 3, 2, 2, 2, 536, 537, 7, 33, 2, 2, 537, 541, 7, 3, 2, 2, 538, 540,
	5, 106, 54, 2, 539, 538, 3, 2, 2, 2, 540, 543, 3, 2, 2, 2, 541, 539, 3,
	2, 2, 2, 541, 542, 3, 2, 2, 2, 542, 544, 3, 2, 2, 2, 543, 541, 3, 2, 2,
	2, 544, 545, 7, 4, 2, 2, 545, 109, 3, 2, 2, 2, 546, 550, 5, 114, 58, 2,
	547, 550, 5, 108, 55, 2, 548, 550, 5, 106, 54, 2, 549, 546, 3, 2, 2, 2,
	549, 547, 3, 2, 2, 2, 549, 548, 3, 2, 2, 2, 550, 111, 3, 2, 2, 2, 551,
	552, 7, 17, 2, 2, 552, 113, 3, 2, 2, 2, 553, 567, 7, 17, 2, 2, 554, 564,
	5, 76, 39, 2, 555, 565, 5, 112, 57, 2, 556, 558, 7, 3, 2, 2, 557, 559,
	5, 110, 56, 2, 558, 557, 3, 2, 2, 2, 559, 560, 3, 2, 2, 2, 560, 558, 3,
	2, 2, 2, 560, 561, 3, 2, 2, 2, 561, 562, 3, 2, 2, 2, 562, 563, 7, 4, 2,
	2, 563, 565, 3, 2, 2, 2, 564, 555, 3, 2, 2, 2, 564, 556, 3, 2, 2, 2, 565,
	567, 3, 2, 2, 2, 566, 553, 3, 2, 2, 2, 566, 554, 3, 2, 2, 2, 567, 115,
	3, 2, 2, 2, 568, 572, 5, 84, 43, 2, 569, 570, 7, 6, 2, 2, 570, 572, 5,
	74, 38, 2, 571, 568, 3, 2, 2, 2, 571, 569, 3, 2, 2, 2, 572, 117, 3, 2,
	2, 2, 573, 574, 5, 116, 59, 2, 574, 575, 5, 26, 14, 2, 575, 119, 3, 2,
	2, 2, 576, 577, 7, 24, 2, 2, 577, 587, 7, 33, 2, 2, 578, 588, 7, 17, 2,
	2, 579, 581, 7, 3, 2, 2, 580, 582, 5, 118, 60, 2, 581, 580, 3, 2, 2, 2,
	582, 583, 3, 2, 2, 2, 583, 581, 3, 2, 2, 2, 583, 584, 3, 2, 2, 2, 584,
	585, 3, 2, 2, 2, 585, 586, 7, 4, 2, 2, 586, 588, 3, 2, 2, 2, 587, 578,
	3, 2, 2, 2, 587, 579, 3, 2, 2, 2, 588, 121, 3, 2, 2, 2, 589, 590, 7, 20,
	2, 2, 590, 592, 7, 67, 2, 2, 591, 593, 5, 26, 14, 2, 592, 591, 3, 2, 2,
	2, 592, 593, 3, 2, 2, 2, 593, 594, 3, 2, 2, 2, 594, 604, 7, 33, 2, 2, 595,
	605, 7, 17, 2, 2, 596, 598, 7, 3, 2, 2, 597, 599, 5, 106, 54, 2, 598, 597,
	3, 2, 2, 2, 599, 600, 3, 2, 2, 2, 600, 598, 3, 2, 2, 2, 600, 601, 3, 2,
	2, 2, 601, 602, 3, 2, 2, 2, 602, 603, 7, 4, 2, 2, 603, 605, 3, 2, 2, 2,
	604, 595, 3, 2, 2, 2, 604, 596, 3, 2, 2, 2, 605, 123, 3, 2, 2, 2, 606,
	614, 7, 3, 2, 2, 607, 615, 5, 46, 24, 2, 608, 615, 5, 62, 32, 2, 609, 615,
	7, 57, 2, 2, 610, 615, 5, 114, 58, 2, 611, 615, 5, 120, 61, 2, 612, 615,
	5, 122, 62, 2, 613, 615, 5, 38, 20, 2, 614, 607, 3, 2, 2, 2, 614, 608,
	3, 2, 2, 2, 614, 609, 3, 2, 2, 2, 614, 610, 3, 2, 2, 2, 614, 611, 3, 2,
	2, 2, 614, 612, 3, 2, 2, 2, 614, 613, 3, 2, 2, 2, 615, 616, 3, 2, 2, 2,
	616, 614, 3, 2, 2, 2, 616, 617, 3, 2, 2, 2, 617, 618, 3, 2, 2, 2, 618,
	619, 7, 4, 2, 2, 619, 125, 3, 2, 2, 2, 620, 622, 7, 57, 2, 2, 621, 620,
	3, 2, 2, 2, 622, 625, 3, 2, 2, 2, 623, 621, 3, 2, 2, 2, 623, 624, 3, 2,
	2, 2, 624, 626, 3, 2, 2, 2, 625, 623, 3, 2, 2, 2, 626, 627, 5, 54, 28,
	2, 627, 628, 7, 33, 2, 2, 628, 629, 5, 124, 63, 2, 629, 127, 3, 2, 2, 2,
	630, 632, 7, 30, 2, 2, 631, 630, 3, 2, 2, 2, 631, 632, 3, 2, 2, 2, 632,
	633, 3, 2, 2, 2, 633, 638, 7, 59, 2, 2, 634, 635, 7, 30, 2, 2, 635, 637,
	7, 59, 2, 2, 636, 634, 3, 2, 2, 2, 637, 640, 3, 2, 2, 2, 638, 636, 3, 2,
	2, 2, 638, 639, 3, 2, 2, 2, 639, 129, 3, 2, 2, 2, 640, 638, 3, 2, 2, 2,
	641, 642, 7, 10, 2, 2, 642, 646, 7, 61, 2, 2, 643, 645, 7, 57, 2, 2, 644,
	643, 3, 2, 2, 2, 645, 648, 3, 2, 2, 2, 646, 644, 3, 2, 2, 2, 646, 647,
	3, 2, 2, 2, 647, 131, 3, 2, 2, 2, 648, 646, 3, 2, 2, 2, 649, 651, 5, 130,
	66, 2, 650, 649, 3, 2, 2, 2, 651, 652, 3, 2, 2, 2, 652, 650, 3, 2, 2, 2,
	652, 653, 3, 2, 2, 2, 653, 133, 3, 2, 2, 2, 654, 656, 5, 132, 67, 2, 655,
	654, 3, 2, 2, 2, 655, 656, 3, 2, 2, 2, 656, 658, 3, 2, 2, 2, 657, 659,
	5, 126, 64, 2, 658, 657, 3, 2, 2, 2, 659, 660, 3, 2, 2, 2, 660, 658, 3,
	2, 2, 2, 660, 661, 3, 2, 2, 2, 661, 662, 3, 2, 2, 2, 662, 663, 7, 2, 2,
	3, 663, 135, 3, 2, 2, 2, 76, 142, 149, 158, 180, 194, 202, 209, 217, 225,
	236, 243, 254, 262, 265, 268, 271, 275, 277, 286, 292, 298, 300, 313, 318,
	321, 330, 334, 341, 346, 355, 369, 377, 389, 394, 396, 399, 403, 408, 411,
	414, 417, 426, 433, 442, 454, 458, 469, 483, 490, 499, 509, 527, 531, 534,
	541, 549, 560, 564, 566, 571, 583, 587, 592, 600, 604, 614, 616, 623, 631,
	638, 646, 652, 655, 660,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "", "", "", "", "'!wrap'", "'!table'", "'!type'", "'import'", "", "",
	"", "", "", "", "'...'", "'set of'", "", "'<->'", "'::'", "'<:'", "'<-'",
	"'.. * <- *'", "'+'", "'~'", "','", "'='", "'$'", "'/'", "'-'", "'*'",
	"':'", "'%'", "'.'", "'!'", "'?'", "'@'", "'&'", "'['", "']'", "'{'", "'}'",
	"'('", "')'", "'<'", "'>'", "'#'", "'|'",
}
var symbolicNames = []string{
	"", "INDENT", "DEDENT", "NativeDataTypes", "HTTP_VERBS", "WRAP", "TABLE",
	"TYPE", "IMPORT", "RETURN", "IF", "ELSE", "FOR", "LOOP", "GROUP", "WHATEVER",
	"SET_OF", "ONE_OF", "DISTANCE", "NAME_SEP", "LESS_COLON", "MEMBER", "COLLECTOR",
	"PLUS", "TILDE", "COMMA", "EQ", "DOLLAR", "FORWARD_SLASH", "MINUS", "STAR",
	"COLON", "PERCENT", "DOT", "EXCLAIM", "QN", "AT", "AMP", "SQ_OPEN", "SQ_CLOSE",
	"CURLY_OPEN", "CURLY_CLOSE", "OPEN_PAREN", "CLOSE_PAREN", "OPEN_ANGLE",
	"CLOSE_ANGLE", "HASH", "PIPE", "DBL_QT", "SINGLE_QT", "EMPTY_LINE", "INDENTED_COMMENT",
	"DIGITS", "QSTRING", "NEWLINE", "SYSL_COMMENT", "TEXT_LINE", "Name", "WS",
	"TEXT", "SKIP_WS", "TEXT_NAME", "POP_WS", "VAR_NAME", "SKIP_WS_2", "EVENT_NAME",
}

var ruleNames = []string{
	"modifier", "size_spec", "modifier_list", "modifiers", "reference", "doc_string",
	"quoted_string", "array_of_strings", "array_of_arrays", "nvp", "attributes",
	"entry", "attribs_or_modifiers", "set_type", "collection_type", "user_defined_type",
	"multi_line_docstring", "annotation_value", "annotation", "annotations",
	"field_type", "field", "table", "package_name", "sub_package", "app_name",
	"name_with_attribs", "model_name", "inplace_table_def", "table_refs", "facade",
	"documentation_stmts", "variable_substitution", "static_path", "query_var",
	"query_param", "http_path", "endpoint_name", "ret_stmt", "target", "target_endpoint",
	"call_stmt", "if_stmt", "if_else", "for_cond", "for_stmt", "http_method_comment",
	"group_stmt", "one_of_case_label", "one_of_cases", "one_of_stmt", "text_stmt",
	"http_statements", "method_def", "endpoint_decl", "shortcut", "api_endpoint",
	"collector_stmt", "collector_stmts", "collector", "event", "app_decl",
	"application", "path", "import_stmt", "imports_decl", "sysl_file",
}
var decisionToDFA = make([]*antlr.DFA, len(deserializedATN.DecisionToState))

func init() {
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

type SyslParser struct {
	*antlr.BaseParser
}

func NewSyslParser(input antlr.TokenStream) *SyslParser {
	this := new(SyslParser)

	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "SyslParser.g4"

	return this
}

// SyslParser tokens.
const (
	SyslParserEOF              = antlr.TokenEOF
	SyslParserINDENT           = 1
	SyslParserDEDENT           = 2
	SyslParserNativeDataTypes  = 3
	SyslParserHTTP_VERBS       = 4
	SyslParserWRAP             = 5
	SyslParserTABLE            = 6
	SyslParserTYPE             = 7
	SyslParserIMPORT           = 8
	SyslParserRETURN           = 9
	SyslParserIF               = 10
	SyslParserELSE             = 11
	SyslParserFOR              = 12
	SyslParserLOOP             = 13
	SyslParserGROUP            = 14
	SyslParserWHATEVER         = 15
	SyslParserSET_OF           = 16
	SyslParserONE_OF           = 17
	SyslParserDISTANCE         = 18
	SyslParserNAME_SEP         = 19
	SyslParserLESS_COLON       = 20
	SyslParserMEMBER           = 21
	SyslParserCOLLECTOR        = 22
	SyslParserPLUS             = 23
	SyslParserTILDE            = 24
	SyslParserCOMMA            = 25
	SyslParserEQ               = 26
	SyslParserDOLLAR           = 27
	SyslParserFORWARD_SLASH    = 28
	SyslParserMINUS            = 29
	SyslParserSTAR             = 30
	SyslParserCOLON            = 31
	SyslParserPERCENT          = 32
	SyslParserDOT              = 33
	SyslParserEXCLAIM          = 34
	SyslParserQN               = 35
	SyslParserAT               = 36
	SyslParserAMP              = 37
	SyslParserSQ_OPEN          = 38
	SyslParserSQ_CLOSE         = 39
	SyslParserCURLY_OPEN       = 40
	SyslParserCURLY_CLOSE      = 41
	SyslParserOPEN_PAREN       = 42
	SyslParserCLOSE_PAREN      = 43
	SyslParserOPEN_ANGLE       = 44
	SyslParserCLOSE_ANGLE      = 45
	SyslParserHASH             = 46
	SyslParserPIPE             = 47
	SyslParserDBL_QT           = 48
	SyslParserSINGLE_QT        = 49
	SyslParserEMPTY_LINE       = 50
	SyslParserINDENTED_COMMENT = 51
	SyslParserDIGITS           = 52
	SyslParserQSTRING          = 53
	SyslParserNEWLINE          = 54
	SyslParserSYSL_COMMENT     = 55
	SyslParserTEXT_LINE        = 56
	SyslParserName             = 57
	SyslParserWS               = 58
	SyslParserTEXT             = 59
	SyslParserSKIP_WS          = 60
	SyslParserTEXT_NAME        = 61
	SyslParserPOP_WS           = 62
	SyslParserVAR_NAME         = 63
	SyslParserSKIP_WS_2        = 64
	SyslParserEVENT_NAME       = 65
)

// SyslParser rules.
const (
	SyslParserRULE_modifier              = 0
	SyslParserRULE_size_spec             = 1
	SyslParserRULE_modifier_list         = 2
	SyslParserRULE_modifiers             = 3
	SyslParserRULE_reference             = 4
	SyslParserRULE_doc_string            = 5
	SyslParserRULE_quoted_string         = 6
	SyslParserRULE_array_of_strings      = 7
	SyslParserRULE_array_of_arrays       = 8
	SyslParserRULE_nvp                   = 9
	SyslParserRULE_attributes            = 10
	SyslParserRULE_entry                 = 11
	SyslParserRULE_attribs_or_modifiers  = 12
	SyslParserRULE_set_type              = 13
	SyslParserRULE_collection_type       = 14
	SyslParserRULE_user_defined_type     = 15
	SyslParserRULE_multi_line_docstring  = 16
	SyslParserRULE_annotation_value      = 17
	SyslParserRULE_annotation            = 18
	SyslParserRULE_annotations           = 19
	SyslParserRULE_field_type            = 20
	SyslParserRULE_field                 = 21
	SyslParserRULE_table                 = 22
	SyslParserRULE_package_name          = 23
	SyslParserRULE_sub_package           = 24
	SyslParserRULE_app_name              = 25
	SyslParserRULE_name_with_attribs     = 26
	SyslParserRULE_model_name            = 27
	SyslParserRULE_inplace_table_def     = 28
	SyslParserRULE_table_refs            = 29
	SyslParserRULE_facade                = 30
	SyslParserRULE_documentation_stmts   = 31
	SyslParserRULE_variable_substitution = 32
	SyslParserRULE_static_path           = 33
	SyslParserRULE_query_var             = 34
	SyslParserRULE_query_param           = 35
	SyslParserRULE_http_path             = 36
	SyslParserRULE_endpoint_name         = 37
	SyslParserRULE_ret_stmt              = 38
	SyslParserRULE_target                = 39
	SyslParserRULE_target_endpoint       = 40
	SyslParserRULE_call_stmt             = 41
	SyslParserRULE_if_stmt               = 42
	SyslParserRULE_if_else               = 43
	SyslParserRULE_for_cond              = 44
	SyslParserRULE_for_stmt              = 45
	SyslParserRULE_http_method_comment   = 46
	SyslParserRULE_group_stmt            = 47
	SyslParserRULE_one_of_case_label     = 48
	SyslParserRULE_one_of_cases          = 49
	SyslParserRULE_one_of_stmt           = 50
	SyslParserRULE_text_stmt             = 51
	SyslParserRULE_http_statements       = 52
	SyslParserRULE_method_def            = 53
	SyslParserRULE_endpoint_decl         = 54
	SyslParserRULE_shortcut              = 55
	SyslParserRULE_api_endpoint          = 56
	SyslParserRULE_collector_stmt        = 57
	SyslParserRULE_collector_stmts       = 58
	SyslParserRULE_collector             = 59
	SyslParserRULE_event                 = 60
	SyslParserRULE_app_decl              = 61
	SyslParserRULE_application           = 62
	SyslParserRULE_path                  = 63
	SyslParserRULE_import_stmt           = 64
	SyslParserRULE_imports_decl          = 65
	SyslParserRULE_sysl_file             = 66
)

// IModifierContext is an interface to support dynamic dispatch.
type IModifierContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsModifierContext differentiates from other interfaces.
	IsModifierContext()
}

type ModifierContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyModifierContext() *ModifierContext {
	var p = new(ModifierContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_modifier
	return p
}

func (*ModifierContext) IsModifierContext() {}

func NewModifierContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ModifierContext {
	var p = new(ModifierContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_modifier

	return p
}

func (s *ModifierContext) GetParser() antlr.Parser { return s.parser }

func (s *ModifierContext) TILDE() antlr.TerminalNode {
	return s.GetToken(SyslParserTILDE, 0)
}

func (s *ModifierContext) AllName() []antlr.TerminalNode {
	return s.GetTokens(SyslParserName)
}

func (s *ModifierContext) Name(i int) antlr.TerminalNode {
	return s.GetToken(SyslParserName, i)
}

func (s *ModifierContext) AllPLUS() []antlr.TerminalNode {
	return s.GetTokens(SyslParserPLUS)
}

func (s *ModifierContext) PLUS(i int) antlr.TerminalNode {
	return s.GetToken(SyslParserPLUS, i)
}

func (s *ModifierContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ModifierContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ModifierContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterModifier(s)
	}
}

func (s *ModifierContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitModifier(s)
	}
}

func (p *SyslParser) Modifier() (localctx IModifierContext) {
	localctx = NewModifierContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, SyslParserRULE_modifier)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(134)
		p.Match(SyslParserTILDE)
	}
	{
		p.SetState(135)
		p.Match(SyslParserName)
	}
	p.SetState(140)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SyslParserPLUS {
		{
			p.SetState(136)
			p.Match(SyslParserPLUS)
		}
		{
			p.SetState(137)
			p.Match(SyslParserName)
		}

		p.SetState(142)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// ISize_specContext is an interface to support dynamic dispatch.
type ISize_specContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSize_specContext differentiates from other interfaces.
	IsSize_specContext()
}

type Size_specContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySize_specContext() *Size_specContext {
	var p = new(Size_specContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_size_spec
	return p
}

func (*Size_specContext) IsSize_specContext() {}

func NewSize_specContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Size_specContext {
	var p = new(Size_specContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_size_spec

	return p
}

func (s *Size_specContext) GetParser() antlr.Parser { return s.parser }

func (s *Size_specContext) OPEN_PAREN() antlr.TerminalNode {
	return s.GetToken(SyslParserOPEN_PAREN, 0)
}

func (s *Size_specContext) AllDIGITS() []antlr.TerminalNode {
	return s.GetTokens(SyslParserDIGITS)
}

func (s *Size_specContext) DIGITS(i int) antlr.TerminalNode {
	return s.GetToken(SyslParserDIGITS, i)
}

func (s *Size_specContext) CLOSE_PAREN() antlr.TerminalNode {
	return s.GetToken(SyslParserCLOSE_PAREN, 0)
}

func (s *Size_specContext) DOT() antlr.TerminalNode {
	return s.GetToken(SyslParserDOT, 0)
}

func (s *Size_specContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Size_specContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Size_specContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterSize_spec(s)
	}
}

func (s *Size_specContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitSize_spec(s)
	}
}

func (p *SyslParser) Size_spec() (localctx ISize_specContext) {
	localctx = NewSize_specContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, SyslParserRULE_size_spec)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(143)
		p.Match(SyslParserOPEN_PAREN)
	}
	{
		p.SetState(144)
		p.Match(SyslParserDIGITS)
	}
	p.SetState(147)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserDOT {
		{
			p.SetState(145)
			p.Match(SyslParserDOT)
		}
		{
			p.SetState(146)
			p.Match(SyslParserDIGITS)
		}

	}
	{
		p.SetState(149)
		p.Match(SyslParserCLOSE_PAREN)
	}

	return localctx
}

// IModifier_listContext is an interface to support dynamic dispatch.
type IModifier_listContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsModifier_listContext differentiates from other interfaces.
	IsModifier_listContext()
}

type Modifier_listContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyModifier_listContext() *Modifier_listContext {
	var p = new(Modifier_listContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_modifier_list
	return p
}

func (*Modifier_listContext) IsModifier_listContext() {}

func NewModifier_listContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Modifier_listContext {
	var p = new(Modifier_listContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_modifier_list

	return p
}

func (s *Modifier_listContext) GetParser() antlr.Parser { return s.parser }

func (s *Modifier_listContext) AllModifier() []IModifierContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IModifierContext)(nil)).Elem())
	var tst = make([]IModifierContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IModifierContext)
		}
	}

	return tst
}

func (s *Modifier_listContext) Modifier(i int) IModifierContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IModifierContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IModifierContext)
}

func (s *Modifier_listContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SyslParserCOMMA)
}

func (s *Modifier_listContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SyslParserCOMMA, i)
}

func (s *Modifier_listContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Modifier_listContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Modifier_listContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterModifier_list(s)
	}
}

func (s *Modifier_listContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitModifier_list(s)
	}
}

func (p *SyslParser) Modifier_list() (localctx IModifier_listContext) {
	localctx = NewModifier_listContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, SyslParserRULE_modifier_list)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(151)
		p.Modifier()
	}
	p.SetState(156)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SyslParserCOMMA {
		{
			p.SetState(152)
			p.Match(SyslParserCOMMA)
		}
		{
			p.SetState(153)
			p.Modifier()
		}

		p.SetState(158)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IModifiersContext is an interface to support dynamic dispatch.
type IModifiersContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsModifiersContext differentiates from other interfaces.
	IsModifiersContext()
}

type ModifiersContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyModifiersContext() *ModifiersContext {
	var p = new(ModifiersContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_modifiers
	return p
}

func (*ModifiersContext) IsModifiersContext() {}

func NewModifiersContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ModifiersContext {
	var p = new(ModifiersContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_modifiers

	return p
}

func (s *ModifiersContext) GetParser() antlr.Parser { return s.parser }

func (s *ModifiersContext) SQ_OPEN() antlr.TerminalNode {
	return s.GetToken(SyslParserSQ_OPEN, 0)
}

func (s *ModifiersContext) Modifier_list() IModifier_listContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IModifier_listContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IModifier_listContext)
}

func (s *ModifiersContext) SQ_CLOSE() antlr.TerminalNode {
	return s.GetToken(SyslParserSQ_CLOSE, 0)
}

func (s *ModifiersContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ModifiersContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ModifiersContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterModifiers(s)
	}
}

func (s *ModifiersContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitModifiers(s)
	}
}

func (p *SyslParser) Modifiers() (localctx IModifiersContext) {
	localctx = NewModifiersContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, SyslParserRULE_modifiers)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(159)
		p.Match(SyslParserSQ_OPEN)
	}
	{
		p.SetState(160)
		p.Modifier_list()
	}
	{
		p.SetState(161)
		p.Match(SyslParserSQ_CLOSE)
	}

	return localctx
}

// IReferenceContext is an interface to support dynamic dispatch.
type IReferenceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetParent_ref returns the parent_ref token.
	GetParent_ref() antlr.Token

	// GetMember returns the member token.
	GetMember() antlr.Token

	// SetParent_ref sets the parent_ref token.
	SetParent_ref(antlr.Token)

	// SetMember sets the member token.
	SetMember(antlr.Token)

	// IsReferenceContext differentiates from other interfaces.
	IsReferenceContext()
}

type ReferenceContext struct {
	*antlr.BaseParserRuleContext
	parser     antlr.Parser
	parent_ref antlr.Token
	member     antlr.Token
}

func NewEmptyReferenceContext() *ReferenceContext {
	var p = new(ReferenceContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_reference
	return p
}

func (*ReferenceContext) IsReferenceContext() {}

func NewReferenceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ReferenceContext {
	var p = new(ReferenceContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_reference

	return p
}

func (s *ReferenceContext) GetParser() antlr.Parser { return s.parser }

func (s *ReferenceContext) GetParent_ref() antlr.Token { return s.parent_ref }

func (s *ReferenceContext) GetMember() antlr.Token { return s.member }

func (s *ReferenceContext) SetParent_ref(v antlr.Token) { s.parent_ref = v }

func (s *ReferenceContext) SetMember(v antlr.Token) { s.member = v }

func (s *ReferenceContext) AllName() []antlr.TerminalNode {
	return s.GetTokens(SyslParserName)
}

func (s *ReferenceContext) Name(i int) antlr.TerminalNode {
	return s.GetToken(SyslParserName, i)
}

func (s *ReferenceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ReferenceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ReferenceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterReference(s)
	}
}

func (s *ReferenceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitReference(s)
	}
}

func (p *SyslParser) Reference() (localctx IReferenceContext) {
	localctx = NewReferenceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, SyslParserRULE_reference)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(163)

		var _m = p.Match(SyslParserName)

		localctx.(*ReferenceContext).parent_ref = _m
	}
	{
		p.SetState(164)
		p.Match(SyslParserDOT)
	}
	{
		p.SetState(165)

		var _m = p.Match(SyslParserName)

		localctx.(*ReferenceContext).member = _m
	}

	return localctx
}

// IDoc_stringContext is an interface to support dynamic dispatch.
type IDoc_stringContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDoc_stringContext differentiates from other interfaces.
	IsDoc_stringContext()
}

type Doc_stringContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDoc_stringContext() *Doc_stringContext {
	var p = new(Doc_stringContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_doc_string
	return p
}

func (*Doc_stringContext) IsDoc_stringContext() {}

func NewDoc_stringContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Doc_stringContext {
	var p = new(Doc_stringContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_doc_string

	return p
}

func (s *Doc_stringContext) GetParser() antlr.Parser { return s.parser }

func (s *Doc_stringContext) PIPE() antlr.TerminalNode {
	return s.GetToken(SyslParserPIPE, 0)
}

func (s *Doc_stringContext) TEXT() antlr.TerminalNode {
	return s.GetToken(SyslParserTEXT, 0)
}

func (s *Doc_stringContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Doc_stringContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Doc_stringContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterDoc_string(s)
	}
}

func (s *Doc_stringContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitDoc_string(s)
	}
}

func (p *SyslParser) Doc_string() (localctx IDoc_stringContext) {
	localctx = NewDoc_stringContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, SyslParserRULE_doc_string)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(167)
		p.Match(SyslParserPIPE)
	}
	{
		p.SetState(168)
		p.Match(SyslParserTEXT)
	}

	return localctx
}

// IQuoted_stringContext is an interface to support dynamic dispatch.
type IQuoted_stringContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsQuoted_stringContext differentiates from other interfaces.
	IsQuoted_stringContext()
}

type Quoted_stringContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQuoted_stringContext() *Quoted_stringContext {
	var p = new(Quoted_stringContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_quoted_string
	return p
}

func (*Quoted_stringContext) IsQuoted_stringContext() {}

func NewQuoted_stringContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Quoted_stringContext {
	var p = new(Quoted_stringContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_quoted_string

	return p
}

func (s *Quoted_stringContext) GetParser() antlr.Parser { return s.parser }

func (s *Quoted_stringContext) QSTRING() antlr.TerminalNode {
	return s.GetToken(SyslParserQSTRING, 0)
}

func (s *Quoted_stringContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Quoted_stringContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Quoted_stringContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterQuoted_string(s)
	}
}

func (s *Quoted_stringContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitQuoted_string(s)
	}
}

func (p *SyslParser) Quoted_string() (localctx IQuoted_stringContext) {
	localctx = NewQuoted_stringContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, SyslParserRULE_quoted_string)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(170)
		p.Match(SyslParserQSTRING)
	}

	return localctx
}

// IArray_of_stringsContext is an interface to support dynamic dispatch.
type IArray_of_stringsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsArray_of_stringsContext differentiates from other interfaces.
	IsArray_of_stringsContext()
}

type Array_of_stringsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArray_of_stringsContext() *Array_of_stringsContext {
	var p = new(Array_of_stringsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_array_of_strings
	return p
}

func (*Array_of_stringsContext) IsArray_of_stringsContext() {}

func NewArray_of_stringsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Array_of_stringsContext {
	var p = new(Array_of_stringsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_array_of_strings

	return p
}

func (s *Array_of_stringsContext) GetParser() antlr.Parser { return s.parser }

func (s *Array_of_stringsContext) SQ_OPEN() antlr.TerminalNode {
	return s.GetToken(SyslParserSQ_OPEN, 0)
}

func (s *Array_of_stringsContext) AllQuoted_string() []IQuoted_stringContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IQuoted_stringContext)(nil)).Elem())
	var tst = make([]IQuoted_stringContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IQuoted_stringContext)
		}
	}

	return tst
}

func (s *Array_of_stringsContext) Quoted_string(i int) IQuoted_stringContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IQuoted_stringContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IQuoted_stringContext)
}

func (s *Array_of_stringsContext) SQ_CLOSE() antlr.TerminalNode {
	return s.GetToken(SyslParserSQ_CLOSE, 0)
}

func (s *Array_of_stringsContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SyslParserCOMMA)
}

func (s *Array_of_stringsContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SyslParserCOMMA, i)
}

func (s *Array_of_stringsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Array_of_stringsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Array_of_stringsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterArray_of_strings(s)
	}
}

func (s *Array_of_stringsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitArray_of_strings(s)
	}
}

func (p *SyslParser) Array_of_strings() (localctx IArray_of_stringsContext) {
	localctx = NewArray_of_stringsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, SyslParserRULE_array_of_strings)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(172)
		p.Match(SyslParserSQ_OPEN)
	}
	{
		p.SetState(173)
		p.Quoted_string()
	}
	p.SetState(178)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SyslParserCOMMA {
		{
			p.SetState(174)
			p.Match(SyslParserCOMMA)
		}
		{
			p.SetState(175)
			p.Quoted_string()
		}

		p.SetState(180)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(181)
		p.Match(SyslParserSQ_CLOSE)
	}

	return localctx
}

// IArray_of_arraysContext is an interface to support dynamic dispatch.
type IArray_of_arraysContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsArray_of_arraysContext differentiates from other interfaces.
	IsArray_of_arraysContext()
}

type Array_of_arraysContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArray_of_arraysContext() *Array_of_arraysContext {
	var p = new(Array_of_arraysContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_array_of_arrays
	return p
}

func (*Array_of_arraysContext) IsArray_of_arraysContext() {}

func NewArray_of_arraysContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Array_of_arraysContext {
	var p = new(Array_of_arraysContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_array_of_arrays

	return p
}

func (s *Array_of_arraysContext) GetParser() antlr.Parser { return s.parser }

func (s *Array_of_arraysContext) SQ_OPEN() antlr.TerminalNode {
	return s.GetToken(SyslParserSQ_OPEN, 0)
}

func (s *Array_of_arraysContext) Array_of_strings() IArray_of_stringsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArray_of_stringsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArray_of_stringsContext)
}

func (s *Array_of_arraysContext) SQ_CLOSE() antlr.TerminalNode {
	return s.GetToken(SyslParserSQ_CLOSE, 0)
}

func (s *Array_of_arraysContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Array_of_arraysContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Array_of_arraysContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterArray_of_arrays(s)
	}
}

func (s *Array_of_arraysContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitArray_of_arrays(s)
	}
}

func (p *SyslParser) Array_of_arrays() (localctx IArray_of_arraysContext) {
	localctx = NewArray_of_arraysContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, SyslParserRULE_array_of_arrays)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(183)
		p.Match(SyslParserSQ_OPEN)
	}
	{
		p.SetState(184)
		p.Array_of_strings()
	}
	{
		p.SetState(185)
		p.Match(SyslParserSQ_CLOSE)
	}

	return localctx
}

// INvpContext is an interface to support dynamic dispatch.
type INvpContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNvpContext differentiates from other interfaces.
	IsNvpContext()
}

type NvpContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNvpContext() *NvpContext {
	var p = new(NvpContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_nvp
	return p
}

func (*NvpContext) IsNvpContext() {}

func NewNvpContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NvpContext {
	var p = new(NvpContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_nvp

	return p
}

func (s *NvpContext) GetParser() antlr.Parser { return s.parser }

func (s *NvpContext) Name() antlr.TerminalNode {
	return s.GetToken(SyslParserName, 0)
}

func (s *NvpContext) EQ() antlr.TerminalNode {
	return s.GetToken(SyslParserEQ, 0)
}

func (s *NvpContext) Quoted_string() IQuoted_stringContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IQuoted_stringContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IQuoted_stringContext)
}

func (s *NvpContext) Array_of_strings() IArray_of_stringsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArray_of_stringsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArray_of_stringsContext)
}

func (s *NvpContext) Array_of_arrays() IArray_of_arraysContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArray_of_arraysContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArray_of_arraysContext)
}

func (s *NvpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NvpContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NvpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterNvp(s)
	}
}

func (s *NvpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitNvp(s)
	}
}

func (p *SyslParser) Nvp() (localctx INvpContext) {
	localctx = NewNvpContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, SyslParserRULE_nvp)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(187)
		p.Match(SyslParserName)
	}
	{
		p.SetState(188)
		p.Match(SyslParserEQ)
	}
	p.SetState(192)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 4, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(189)
			p.Quoted_string()
		}

	case 2:
		{
			p.SetState(190)
			p.Array_of_strings()
		}

	case 3:
		{
			p.SetState(191)
			p.Array_of_arrays()
		}

	}

	return localctx
}

// IAttributesContext is an interface to support dynamic dispatch.
type IAttributesContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAttributesContext differentiates from other interfaces.
	IsAttributesContext()
}

type AttributesContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAttributesContext() *AttributesContext {
	var p = new(AttributesContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_attributes
	return p
}

func (*AttributesContext) IsAttributesContext() {}

func NewAttributesContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AttributesContext {
	var p = new(AttributesContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_attributes

	return p
}

func (s *AttributesContext) GetParser() antlr.Parser { return s.parser }

func (s *AttributesContext) SQ_OPEN() antlr.TerminalNode {
	return s.GetToken(SyslParserSQ_OPEN, 0)
}

func (s *AttributesContext) AllNvp() []INvpContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*INvpContext)(nil)).Elem())
	var tst = make([]INvpContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(INvpContext)
		}
	}

	return tst
}

func (s *AttributesContext) Nvp(i int) INvpContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INvpContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(INvpContext)
}

func (s *AttributesContext) SQ_CLOSE() antlr.TerminalNode {
	return s.GetToken(SyslParserSQ_CLOSE, 0)
}

func (s *AttributesContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SyslParserCOMMA)
}

func (s *AttributesContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SyslParserCOMMA, i)
}

func (s *AttributesContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AttributesContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AttributesContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterAttributes(s)
	}
}

func (s *AttributesContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitAttributes(s)
	}
}

func (p *SyslParser) Attributes() (localctx IAttributesContext) {
	localctx = NewAttributesContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, SyslParserRULE_attributes)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(194)
		p.Match(SyslParserSQ_OPEN)
	}
	{
		p.SetState(195)
		p.Nvp()
	}
	p.SetState(200)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SyslParserCOMMA {
		{
			p.SetState(196)
			p.Match(SyslParserCOMMA)
		}
		{
			p.SetState(197)
			p.Nvp()
		}

		p.SetState(202)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(203)
		p.Match(SyslParserSQ_CLOSE)
	}

	return localctx
}

// IEntryContext is an interface to support dynamic dispatch.
type IEntryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEntryContext differentiates from other interfaces.
	IsEntryContext()
}

type EntryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEntryContext() *EntryContext {
	var p = new(EntryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_entry
	return p
}

func (*EntryContext) IsEntryContext() {}

func NewEntryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EntryContext {
	var p = new(EntryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_entry

	return p
}

func (s *EntryContext) GetParser() antlr.Parser { return s.parser }

func (s *EntryContext) Nvp() INvpContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INvpContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INvpContext)
}

func (s *EntryContext) Modifier() IModifierContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IModifierContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IModifierContext)
}

func (s *EntryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EntryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EntryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterEntry(s)
	}
}

func (s *EntryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitEntry(s)
	}
}

func (p *SyslParser) Entry() (localctx IEntryContext) {
	localctx = NewEntryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, SyslParserRULE_entry)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(207)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SyslParserName:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(205)
			p.Nvp()
		}

	case SyslParserTILDE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(206)
			p.Modifier()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IAttribs_or_modifiersContext is an interface to support dynamic dispatch.
type IAttribs_or_modifiersContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAttribs_or_modifiersContext differentiates from other interfaces.
	IsAttribs_or_modifiersContext()
}

type Attribs_or_modifiersContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAttribs_or_modifiersContext() *Attribs_or_modifiersContext {
	var p = new(Attribs_or_modifiersContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_attribs_or_modifiers
	return p
}

func (*Attribs_or_modifiersContext) IsAttribs_or_modifiersContext() {}

func NewAttribs_or_modifiersContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Attribs_or_modifiersContext {
	var p = new(Attribs_or_modifiersContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_attribs_or_modifiers

	return p
}

func (s *Attribs_or_modifiersContext) GetParser() antlr.Parser { return s.parser }

func (s *Attribs_or_modifiersContext) SQ_OPEN() antlr.TerminalNode {
	return s.GetToken(SyslParserSQ_OPEN, 0)
}

func (s *Attribs_or_modifiersContext) AllEntry() []IEntryContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IEntryContext)(nil)).Elem())
	var tst = make([]IEntryContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IEntryContext)
		}
	}

	return tst
}

func (s *Attribs_or_modifiersContext) Entry(i int) IEntryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEntryContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IEntryContext)
}

func (s *Attribs_or_modifiersContext) SQ_CLOSE() antlr.TerminalNode {
	return s.GetToken(SyslParserSQ_CLOSE, 0)
}

func (s *Attribs_or_modifiersContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SyslParserCOMMA)
}

func (s *Attribs_or_modifiersContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SyslParserCOMMA, i)
}

func (s *Attribs_or_modifiersContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Attribs_or_modifiersContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Attribs_or_modifiersContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterAttribs_or_modifiers(s)
	}
}

func (s *Attribs_or_modifiersContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitAttribs_or_modifiers(s)
	}
}

func (p *SyslParser) Attribs_or_modifiers() (localctx IAttribs_or_modifiersContext) {
	localctx = NewAttribs_or_modifiersContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, SyslParserRULE_attribs_or_modifiers)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(209)
		p.Match(SyslParserSQ_OPEN)
	}
	{
		p.SetState(210)
		p.Entry()
	}
	p.SetState(215)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SyslParserCOMMA {
		{
			p.SetState(211)
			p.Match(SyslParserCOMMA)
		}
		{
			p.SetState(212)
			p.Entry()
		}

		p.SetState(217)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(218)
		p.Match(SyslParserSQ_CLOSE)
	}

	return localctx
}

// ISet_typeContext is an interface to support dynamic dispatch.
type ISet_typeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSet_typeContext differentiates from other interfaces.
	IsSet_typeContext()
}

type Set_typeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySet_typeContext() *Set_typeContext {
	var p = new(Set_typeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_set_type
	return p
}

func (*Set_typeContext) IsSet_typeContext() {}

func NewSet_typeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Set_typeContext {
	var p = new(Set_typeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_set_type

	return p
}

func (s *Set_typeContext) GetParser() antlr.Parser { return s.parser }

func (s *Set_typeContext) SET_OF() antlr.TerminalNode {
	return s.GetToken(SyslParserSET_OF, 0)
}

func (s *Set_typeContext) Name() antlr.TerminalNode {
	return s.GetToken(SyslParserName, 0)
}

func (s *Set_typeContext) NativeDataTypes() antlr.TerminalNode {
	return s.GetToken(SyslParserNativeDataTypes, 0)
}

func (s *Set_typeContext) Size_spec() ISize_specContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISize_specContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISize_specContext)
}

func (s *Set_typeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Set_typeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Set_typeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterSet_type(s)
	}
}

func (s *Set_typeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitSet_type(s)
	}
}

func (p *SyslParser) Set_type() (localctx ISet_typeContext) {
	localctx = NewSet_typeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, SyslParserRULE_set_type)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(220)
		p.Match(SyslParserSET_OF)
	}
	p.SetState(221)
	_la = p.GetTokenStream().LA(1)

	if !(_la == SyslParserNativeDataTypes || _la == SyslParserName) {
		p.GetErrorHandler().RecoverInline(p)
	} else {
		p.GetErrorHandler().ReportMatch(p)
		p.Consume()
	}
	p.SetState(223)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserOPEN_PAREN {
		{
			p.SetState(222)
			p.Size_spec()
		}

	}

	return localctx
}

// ICollection_typeContext is an interface to support dynamic dispatch.
type ICollection_typeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCollection_typeContext differentiates from other interfaces.
	IsCollection_typeContext()
}

type Collection_typeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCollection_typeContext() *Collection_typeContext {
	var p = new(Collection_typeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_collection_type
	return p
}

func (*Collection_typeContext) IsCollection_typeContext() {}

func NewCollection_typeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Collection_typeContext {
	var p = new(Collection_typeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_collection_type

	return p
}

func (s *Collection_typeContext) GetParser() antlr.Parser { return s.parser }

func (s *Collection_typeContext) Set_type() ISet_typeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISet_typeContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISet_typeContext)
}

func (s *Collection_typeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Collection_typeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Collection_typeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterCollection_type(s)
	}
}

func (s *Collection_typeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitCollection_type(s)
	}
}

func (p *SyslParser) Collection_type() (localctx ICollection_typeContext) {
	localctx = NewCollection_typeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, SyslParserRULE_collection_type)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(225)
		p.Set_type()
	}

	return localctx
}

// IUser_defined_typeContext is an interface to support dynamic dispatch.
type IUser_defined_typeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsUser_defined_typeContext differentiates from other interfaces.
	IsUser_defined_typeContext()
}

type User_defined_typeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUser_defined_typeContext() *User_defined_typeContext {
	var p = new(User_defined_typeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_user_defined_type
	return p
}

func (*User_defined_typeContext) IsUser_defined_typeContext() {}

func NewUser_defined_typeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *User_defined_typeContext {
	var p = new(User_defined_typeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_user_defined_type

	return p
}

func (s *User_defined_typeContext) GetParser() antlr.Parser { return s.parser }

func (s *User_defined_typeContext) Name() antlr.TerminalNode {
	return s.GetToken(SyslParserName, 0)
}

func (s *User_defined_typeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *User_defined_typeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *User_defined_typeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterUser_defined_type(s)
	}
}

func (s *User_defined_typeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitUser_defined_type(s)
	}
}

func (p *SyslParser) User_defined_type() (localctx IUser_defined_typeContext) {
	localctx = NewUser_defined_typeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, SyslParserRULE_user_defined_type)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(227)
		p.Match(SyslParserName)
	}

	return localctx
}

// IMulti_line_docstringContext is an interface to support dynamic dispatch.
type IMulti_line_docstringContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMulti_line_docstringContext differentiates from other interfaces.
	IsMulti_line_docstringContext()
}

type Multi_line_docstringContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMulti_line_docstringContext() *Multi_line_docstringContext {
	var p = new(Multi_line_docstringContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_multi_line_docstring
	return p
}

func (*Multi_line_docstringContext) IsMulti_line_docstringContext() {}

func NewMulti_line_docstringContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Multi_line_docstringContext {
	var p = new(Multi_line_docstringContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_multi_line_docstring

	return p
}

func (s *Multi_line_docstringContext) GetParser() antlr.Parser { return s.parser }

func (s *Multi_line_docstringContext) COLON() antlr.TerminalNode {
	return s.GetToken(SyslParserCOLON, 0)
}

func (s *Multi_line_docstringContext) INDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserINDENT, 0)
}

func (s *Multi_line_docstringContext) DEDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserDEDENT, 0)
}

func (s *Multi_line_docstringContext) AllDoc_string() []IDoc_stringContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IDoc_stringContext)(nil)).Elem())
	var tst = make([]IDoc_stringContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IDoc_stringContext)
		}
	}

	return tst
}

func (s *Multi_line_docstringContext) Doc_string(i int) IDoc_stringContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDoc_stringContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IDoc_stringContext)
}

func (s *Multi_line_docstringContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Multi_line_docstringContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Multi_line_docstringContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterMulti_line_docstring(s)
	}
}

func (s *Multi_line_docstringContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitMulti_line_docstring(s)
	}
}

func (p *SyslParser) Multi_line_docstring() (localctx IMulti_line_docstringContext) {
	localctx = NewMulti_line_docstringContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, SyslParserRULE_multi_line_docstring)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(229)
		p.Match(SyslParserCOLON)
	}
	{
		p.SetState(230)
		p.Match(SyslParserINDENT)
	}
	p.SetState(232)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == SyslParserPIPE {
		{
			p.SetState(231)
			p.Doc_string()
		}

		p.SetState(234)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(236)
		p.Match(SyslParserDEDENT)
	}

	return localctx
}

// IAnnotation_valueContext is an interface to support dynamic dispatch.
type IAnnotation_valueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAnnotation_valueContext differentiates from other interfaces.
	IsAnnotation_valueContext()
}

type Annotation_valueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAnnotation_valueContext() *Annotation_valueContext {
	var p = new(Annotation_valueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_annotation_value
	return p
}

func (*Annotation_valueContext) IsAnnotation_valueContext() {}

func NewAnnotation_valueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Annotation_valueContext {
	var p = new(Annotation_valueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_annotation_value

	return p
}

func (s *Annotation_valueContext) GetParser() antlr.Parser { return s.parser }

func (s *Annotation_valueContext) QSTRING() antlr.TerminalNode {
	return s.GetToken(SyslParserQSTRING, 0)
}

func (s *Annotation_valueContext) Array_of_strings() IArray_of_stringsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArray_of_stringsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArray_of_stringsContext)
}

func (s *Annotation_valueContext) Multi_line_docstring() IMulti_line_docstringContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMulti_line_docstringContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMulti_line_docstringContext)
}

func (s *Annotation_valueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Annotation_valueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Annotation_valueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterAnnotation_value(s)
	}
}

func (s *Annotation_valueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitAnnotation_value(s)
	}
}

func (p *SyslParser) Annotation_value() (localctx IAnnotation_valueContext) {
	localctx = NewAnnotation_valueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, SyslParserRULE_annotation_value)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(241)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SyslParserQSTRING:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(238)
			p.Match(SyslParserQSTRING)
		}

	case SyslParserSQ_OPEN:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(239)
			p.Array_of_strings()
		}

	case SyslParserCOLON:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(240)
			p.Multi_line_docstring()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IAnnotationContext is an interface to support dynamic dispatch.
type IAnnotationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAnnotationContext differentiates from other interfaces.
	IsAnnotationContext()
}

type AnnotationContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAnnotationContext() *AnnotationContext {
	var p = new(AnnotationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_annotation
	return p
}

func (*AnnotationContext) IsAnnotationContext() {}

func NewAnnotationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AnnotationContext {
	var p = new(AnnotationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_annotation

	return p
}

func (s *AnnotationContext) GetParser() antlr.Parser { return s.parser }

func (s *AnnotationContext) AT() antlr.TerminalNode {
	return s.GetToken(SyslParserAT, 0)
}

func (s *AnnotationContext) VAR_NAME() antlr.TerminalNode {
	return s.GetToken(SyslParserVAR_NAME, 0)
}

func (s *AnnotationContext) EQ() antlr.TerminalNode {
	return s.GetToken(SyslParserEQ, 0)
}

func (s *AnnotationContext) Annotation_value() IAnnotation_valueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAnnotation_valueContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAnnotation_valueContext)
}

func (s *AnnotationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AnnotationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AnnotationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterAnnotation(s)
	}
}

func (s *AnnotationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitAnnotation(s)
	}
}

func (p *SyslParser) Annotation() (localctx IAnnotationContext) {
	localctx = NewAnnotationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, SyslParserRULE_annotation)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(243)
		p.Match(SyslParserAT)
	}
	{
		p.SetState(244)
		p.Match(SyslParserVAR_NAME)
	}
	{
		p.SetState(245)
		p.Match(SyslParserEQ)
	}
	{
		p.SetState(246)
		p.Annotation_value()
	}

	return localctx
}

// IAnnotationsContext is an interface to support dynamic dispatch.
type IAnnotationsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAnnotationsContext differentiates from other interfaces.
	IsAnnotationsContext()
}

type AnnotationsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAnnotationsContext() *AnnotationsContext {
	var p = new(AnnotationsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_annotations
	return p
}

func (*AnnotationsContext) IsAnnotationsContext() {}

func NewAnnotationsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AnnotationsContext {
	var p = new(AnnotationsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_annotations

	return p
}

func (s *AnnotationsContext) GetParser() antlr.Parser { return s.parser }

func (s *AnnotationsContext) INDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserINDENT, 0)
}

func (s *AnnotationsContext) DEDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserDEDENT, 0)
}

func (s *AnnotationsContext) AllAnnotation() []IAnnotationContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IAnnotationContext)(nil)).Elem())
	var tst = make([]IAnnotationContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IAnnotationContext)
		}
	}

	return tst
}

func (s *AnnotationsContext) Annotation(i int) IAnnotationContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAnnotationContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IAnnotationContext)
}

func (s *AnnotationsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AnnotationsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AnnotationsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterAnnotations(s)
	}
}

func (s *AnnotationsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitAnnotations(s)
	}
}

func (p *SyslParser) Annotations() (localctx IAnnotationsContext) {
	localctx = NewAnnotationsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, SyslParserRULE_annotations)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(248)
		p.Match(SyslParserINDENT)
	}
	p.SetState(250)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == SyslParserAT {
		{
			p.SetState(249)
			p.Annotation()
		}

		p.SetState(252)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(254)
		p.Match(SyslParserDEDENT)
	}

	return localctx
}

// IField_typeContext is an interface to support dynamic dispatch.
type IField_typeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsField_typeContext differentiates from other interfaces.
	IsField_typeContext()
}

type Field_typeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyField_typeContext() *Field_typeContext {
	var p = new(Field_typeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_field_type
	return p
}

func (*Field_typeContext) IsField_typeContext() {}

func NewField_typeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Field_typeContext {
	var p = new(Field_typeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_field_type

	return p
}

func (s *Field_typeContext) GetParser() antlr.Parser { return s.parser }

func (s *Field_typeContext) Collection_type() ICollection_typeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICollection_typeContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICollection_typeContext)
}

func (s *Field_typeContext) Reference() IReferenceContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IReferenceContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IReferenceContext)
}

func (s *Field_typeContext) NativeDataTypes() antlr.TerminalNode {
	return s.GetToken(SyslParserNativeDataTypes, 0)
}

func (s *Field_typeContext) User_defined_type() IUser_defined_typeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IUser_defined_typeContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IUser_defined_typeContext)
}

func (s *Field_typeContext) Size_spec() ISize_specContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISize_specContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISize_specContext)
}

func (s *Field_typeContext) QN() antlr.TerminalNode {
	return s.GetToken(SyslParserQN, 0)
}

func (s *Field_typeContext) Attribs_or_modifiers() IAttribs_or_modifiersContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAttribs_or_modifiersContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAttribs_or_modifiersContext)
}

func (s *Field_typeContext) COLON() antlr.TerminalNode {
	return s.GetToken(SyslParserCOLON, 0)
}

func (s *Field_typeContext) Annotations() IAnnotationsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAnnotationsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAnnotationsContext)
}

func (s *Field_typeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Field_typeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Field_typeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterField_type(s)
	}
}

func (s *Field_typeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitField_type(s)
	}
}

func (p *SyslParser) Field_type() (localctx IField_typeContext) {
	localctx = NewField_typeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, SyslParserRULE_field_type)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(275)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SyslParserSET_OF:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(256)
			p.Collection_type()
		}

	case SyslParserNativeDataTypes, SyslParserName:
		p.EnterOuterAlt(localctx, 2)
		p.SetState(260)
		p.GetErrorHandler().Sync(p)
		switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 12, p.GetParserRuleContext()) {
		case 1:
			{
				p.SetState(257)
				p.Reference()
			}

		case 2:
			{
				p.SetState(258)
				p.Match(SyslParserNativeDataTypes)
			}

		case 3:
			{
				p.SetState(259)
				p.User_defined_type()
			}

		}
		p.SetState(263)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SyslParserOPEN_PAREN {
			{
				p.SetState(262)
				p.Size_spec()
			}

		}
		p.SetState(266)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SyslParserQN {
			{
				p.SetState(265)
				p.Match(SyslParserQN)
			}

		}
		p.SetState(269)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SyslParserSQ_OPEN {
			{
				p.SetState(268)
				p.Attribs_or_modifiers()
			}

		}
		p.SetState(273)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SyslParserCOLON {
			{
				p.SetState(271)
				p.Match(SyslParserCOLON)
			}
			{
				p.SetState(272)
				p.Annotations()
			}

		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IFieldContext is an interface to support dynamic dispatch.
type IFieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFieldContext differentiates from other interfaces.
	IsFieldContext()
}

type FieldContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldContext() *FieldContext {
	var p = new(FieldContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_field
	return p
}

func (*FieldContext) IsFieldContext() {}

func NewFieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldContext {
	var p = new(FieldContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_field

	return p
}

func (s *FieldContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldContext) Name() antlr.TerminalNode {
	return s.GetToken(SyslParserName, 0)
}

func (s *FieldContext) LESS_COLON() antlr.TerminalNode {
	return s.GetToken(SyslParserLESS_COLON, 0)
}

func (s *FieldContext) Field_type() IField_typeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IField_typeContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IField_typeContext)
}

func (s *FieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterField(s)
	}
}

func (s *FieldContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitField(s)
	}
}

func (p *SyslParser) Field() (localctx IFieldContext) {
	localctx = NewFieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, SyslParserRULE_field)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(277)
		p.Match(SyslParserName)
	}
	{
		p.SetState(278)
		p.Match(SyslParserLESS_COLON)
	}
	{
		p.SetState(279)
		p.Field_type()
	}

	return localctx
}

// ITableContext is an interface to support dynamic dispatch.
type ITableContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTableContext differentiates from other interfaces.
	IsTableContext()
}

type TableContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTableContext() *TableContext {
	var p = new(TableContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_table
	return p
}

func (*TableContext) IsTableContext() {}

func NewTableContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TableContext {
	var p = new(TableContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_table

	return p
}

func (s *TableContext) GetParser() antlr.Parser { return s.parser }

func (s *TableContext) Name() antlr.TerminalNode {
	return s.GetToken(SyslParserName, 0)
}

func (s *TableContext) COLON() antlr.TerminalNode {
	return s.GetToken(SyslParserCOLON, 0)
}

func (s *TableContext) INDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserINDENT, 0)
}

func (s *TableContext) DEDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserDEDENT, 0)
}

func (s *TableContext) TABLE() antlr.TerminalNode {
	return s.GetToken(SyslParserTABLE, 0)
}

func (s *TableContext) TYPE() antlr.TerminalNode {
	return s.GetToken(SyslParserTYPE, 0)
}

func (s *TableContext) AllSYSL_COMMENT() []antlr.TerminalNode {
	return s.GetTokens(SyslParserSYSL_COMMENT)
}

func (s *TableContext) SYSL_COMMENT(i int) antlr.TerminalNode {
	return s.GetToken(SyslParserSYSL_COMMENT, i)
}

func (s *TableContext) Attribs_or_modifiers() IAttribs_or_modifiersContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAttribs_or_modifiersContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAttribs_or_modifiersContext)
}

func (s *TableContext) AllField() []IFieldContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IFieldContext)(nil)).Elem())
	var tst = make([]IFieldContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IFieldContext)
		}
	}

	return tst
}

func (s *TableContext) Field(i int) IFieldContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFieldContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IFieldContext)
}

func (s *TableContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TableContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TableContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterTable(s)
	}
}

func (s *TableContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitTable(s)
	}
}

func (p *SyslParser) Table() (localctx ITableContext) {
	localctx = NewTableContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, SyslParserRULE_table)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(284)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SyslParserSYSL_COMMENT {
		{
			p.SetState(281)
			p.Match(SyslParserSYSL_COMMENT)
		}

		p.SetState(286)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(287)
	_la = p.GetTokenStream().LA(1)

	if !(_la == SyslParserTABLE || _la == SyslParserTYPE) {
		p.GetErrorHandler().RecoverInline(p)
	} else {
		p.GetErrorHandler().ReportMatch(p)
		p.Consume()
	}
	{
		p.SetState(288)
		p.Match(SyslParserName)
	}
	p.SetState(290)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserSQ_OPEN {
		{
			p.SetState(289)
			p.Attribs_or_modifiers()
		}

	}
	{
		p.SetState(292)
		p.Match(SyslParserCOLON)
	}
	{
		p.SetState(293)
		p.Match(SyslParserINDENT)
	}
	p.SetState(296)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == SyslParserSYSL_COMMENT || _la == SyslParserName {
		p.SetState(296)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case SyslParserSYSL_COMMENT:
			{
				p.SetState(294)
				p.Match(SyslParserSYSL_COMMENT)
			}

		case SyslParserName:
			{
				p.SetState(295)
				p.Field()
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(298)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(300)
		p.Match(SyslParserDEDENT)
	}

	return localctx
}

// IPackage_nameContext is an interface to support dynamic dispatch.
type IPackage_nameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPackage_nameContext differentiates from other interfaces.
	IsPackage_nameContext()
}

type Package_nameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPackage_nameContext() *Package_nameContext {
	var p = new(Package_nameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_package_name
	return p
}

func (*Package_nameContext) IsPackage_nameContext() {}

func NewPackage_nameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Package_nameContext {
	var p = new(Package_nameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_package_name

	return p
}

func (s *Package_nameContext) GetParser() antlr.Parser { return s.parser }

func (s *Package_nameContext) Name() antlr.TerminalNode {
	return s.GetToken(SyslParserName, 0)
}

func (s *Package_nameContext) TEXT_LINE() antlr.TerminalNode {
	return s.GetToken(SyslParserTEXT_LINE, 0)
}

func (s *Package_nameContext) TEXT_NAME() antlr.TerminalNode {
	return s.GetToken(SyslParserTEXT_NAME, 0)
}

func (s *Package_nameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Package_nameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Package_nameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterPackage_name(s)
	}
}

func (s *Package_nameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitPackage_name(s)
	}
}

func (p *SyslParser) Package_name() (localctx IPackage_nameContext) {
	localctx = NewPackage_nameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, SyslParserRULE_package_name)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(302)
	_la = p.GetTokenStream().LA(1)

	if !(((_la-56)&-(0x1f+1)) == 0 && ((1<<uint((_la-56)))&((1<<(SyslParserTEXT_LINE-56))|(1<<(SyslParserName-56))|(1<<(SyslParserTEXT_NAME-56)))) != 0) {
		p.GetErrorHandler().RecoverInline(p)
	} else {
		p.GetErrorHandler().ReportMatch(p)
		p.Consume()
	}

	return localctx
}

// ISub_packageContext is an interface to support dynamic dispatch.
type ISub_packageContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSub_packageContext differentiates from other interfaces.
	IsSub_packageContext()
}

type Sub_packageContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySub_packageContext() *Sub_packageContext {
	var p = new(Sub_packageContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_sub_package
	return p
}

func (*Sub_packageContext) IsSub_packageContext() {}

func NewSub_packageContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Sub_packageContext {
	var p = new(Sub_packageContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_sub_package

	return p
}

func (s *Sub_packageContext) GetParser() antlr.Parser { return s.parser }

func (s *Sub_packageContext) NAME_SEP() antlr.TerminalNode {
	return s.GetToken(SyslParserNAME_SEP, 0)
}

func (s *Sub_packageContext) Package_name() IPackage_nameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPackage_nameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPackage_nameContext)
}

func (s *Sub_packageContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Sub_packageContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Sub_packageContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterSub_package(s)
	}
}

func (s *Sub_packageContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitSub_package(s)
	}
}

func (p *SyslParser) Sub_package() (localctx ISub_packageContext) {
	localctx = NewSub_packageContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, SyslParserRULE_sub_package)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(304)
		p.Match(SyslParserNAME_SEP)
	}
	{
		p.SetState(305)
		p.Package_name()
	}

	return localctx
}

// IApp_nameContext is an interface to support dynamic dispatch.
type IApp_nameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsApp_nameContext differentiates from other interfaces.
	IsApp_nameContext()
}

type App_nameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyApp_nameContext() *App_nameContext {
	var p = new(App_nameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_app_name
	return p
}

func (*App_nameContext) IsApp_nameContext() {}

func NewApp_nameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *App_nameContext {
	var p = new(App_nameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_app_name

	return p
}

func (s *App_nameContext) GetParser() antlr.Parser { return s.parser }

func (s *App_nameContext) Package_name() IPackage_nameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPackage_nameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPackage_nameContext)
}

func (s *App_nameContext) AllSub_package() []ISub_packageContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ISub_packageContext)(nil)).Elem())
	var tst = make([]ISub_packageContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ISub_packageContext)
		}
	}

	return tst
}

func (s *App_nameContext) Sub_package(i int) ISub_packageContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISub_packageContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ISub_packageContext)
}

func (s *App_nameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *App_nameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *App_nameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterApp_name(s)
	}
}

func (s *App_nameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitApp_name(s)
	}
}

func (p *SyslParser) App_name() (localctx IApp_nameContext) {
	localctx = NewApp_nameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, SyslParserRULE_app_name)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(307)
		p.Package_name()
	}
	p.SetState(311)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SyslParserNAME_SEP {
		{
			p.SetState(308)
			p.Sub_package()
		}

		p.SetState(313)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IName_with_attribsContext is an interface to support dynamic dispatch.
type IName_with_attribsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsName_with_attribsContext differentiates from other interfaces.
	IsName_with_attribsContext()
}

type Name_with_attribsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyName_with_attribsContext() *Name_with_attribsContext {
	var p = new(Name_with_attribsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_name_with_attribs
	return p
}

func (*Name_with_attribsContext) IsName_with_attribsContext() {}

func NewName_with_attribsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Name_with_attribsContext {
	var p = new(Name_with_attribsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_name_with_attribs

	return p
}

func (s *Name_with_attribsContext) GetParser() antlr.Parser { return s.parser }

func (s *Name_with_attribsContext) App_name() IApp_nameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IApp_nameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IApp_nameContext)
}

func (s *Name_with_attribsContext) QSTRING() antlr.TerminalNode {
	return s.GetToken(SyslParserQSTRING, 0)
}

func (s *Name_with_attribsContext) Attribs_or_modifiers() IAttribs_or_modifiersContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAttribs_or_modifiersContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAttribs_or_modifiersContext)
}

func (s *Name_with_attribsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Name_with_attribsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Name_with_attribsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterName_with_attribs(s)
	}
}

func (s *Name_with_attribsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitName_with_attribs(s)
	}
}

func (p *SyslParser) Name_with_attribs() (localctx IName_with_attribsContext) {
	localctx = NewName_with_attribsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, SyslParserRULE_name_with_attribs)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(314)
		p.App_name()
	}
	p.SetState(316)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserQSTRING {
		{
			p.SetState(315)
			p.Match(SyslParserQSTRING)
		}

	}
	p.SetState(319)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserSQ_OPEN {
		{
			p.SetState(318)
			p.Attribs_or_modifiers()
		}

	}

	return localctx
}

// IModel_nameContext is an interface to support dynamic dispatch.
type IModel_nameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsModel_nameContext differentiates from other interfaces.
	IsModel_nameContext()
}

type Model_nameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyModel_nameContext() *Model_nameContext {
	var p = new(Model_nameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_model_name
	return p
}

func (*Model_nameContext) IsModel_nameContext() {}

func NewModel_nameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Model_nameContext {
	var p = new(Model_nameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_model_name

	return p
}

func (s *Model_nameContext) GetParser() antlr.Parser { return s.parser }

func (s *Model_nameContext) Name() antlr.TerminalNode {
	return s.GetToken(SyslParserName, 0)
}

func (s *Model_nameContext) COLON() antlr.TerminalNode {
	return s.GetToken(SyslParserCOLON, 0)
}

func (s *Model_nameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Model_nameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Model_nameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterModel_name(s)
	}
}

func (s *Model_nameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitModel_name(s)
	}
}

func (p *SyslParser) Model_name() (localctx IModel_nameContext) {
	localctx = NewModel_nameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, SyslParserRULE_model_name)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(321)
		p.Match(SyslParserName)
	}
	{
		p.SetState(322)
		p.Match(SyslParserCOLON)
	}

	return localctx
}

// IInplace_table_defContext is an interface to support dynamic dispatch.
type IInplace_table_defContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsInplace_table_defContext differentiates from other interfaces.
	IsInplace_table_defContext()
}

type Inplace_table_defContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyInplace_table_defContext() *Inplace_table_defContext {
	var p = new(Inplace_table_defContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_inplace_table_def
	return p
}

func (*Inplace_table_defContext) IsInplace_table_defContext() {}

func NewInplace_table_defContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Inplace_table_defContext {
	var p = new(Inplace_table_defContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_inplace_table_def

	return p
}

func (s *Inplace_table_defContext) GetParser() antlr.Parser { return s.parser }

func (s *Inplace_table_defContext) COLON() antlr.TerminalNode {
	return s.GetToken(SyslParserCOLON, 0)
}

func (s *Inplace_table_defContext) INDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserINDENT, 0)
}

func (s *Inplace_table_defContext) DEDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserDEDENT, 0)
}

func (s *Inplace_table_defContext) AllName() []antlr.TerminalNode {
	return s.GetTokens(SyslParserName)
}

func (s *Inplace_table_defContext) Name(i int) antlr.TerminalNode {
	return s.GetToken(SyslParserName, i)
}

func (s *Inplace_table_defContext) AllAttribs_or_modifiers() []IAttribs_or_modifiersContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IAttribs_or_modifiersContext)(nil)).Elem())
	var tst = make([]IAttribs_or_modifiersContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IAttribs_or_modifiersContext)
		}
	}

	return tst
}

func (s *Inplace_table_defContext) Attribs_or_modifiers(i int) IAttribs_or_modifiersContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAttribs_or_modifiersContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IAttribs_or_modifiersContext)
}

func (s *Inplace_table_defContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Inplace_table_defContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Inplace_table_defContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterInplace_table_def(s)
	}
}

func (s *Inplace_table_defContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitInplace_table_def(s)
	}
}

func (p *SyslParser) Inplace_table_def() (localctx IInplace_table_defContext) {
	localctx = NewInplace_table_defContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, SyslParserRULE_inplace_table_def)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(324)
		p.Match(SyslParserCOLON)
	}
	{
		p.SetState(325)
		p.Match(SyslParserINDENT)
	}
	p.SetState(330)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == SyslParserName {
		{
			p.SetState(326)
			p.Match(SyslParserName)
		}
		p.SetState(328)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SyslParserSQ_OPEN {
			{
				p.SetState(327)
				p.Attribs_or_modifiers()
			}

		}

		p.SetState(332)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(334)
		p.Match(SyslParserDEDENT)
	}

	return localctx
}

// ITable_refsContext is an interface to support dynamic dispatch.
type ITable_refsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTable_refsContext differentiates from other interfaces.
	IsTable_refsContext()
}

type Table_refsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTable_refsContext() *Table_refsContext {
	var p = new(Table_refsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_table_refs
	return p
}

func (*Table_refsContext) IsTable_refsContext() {}

func NewTable_refsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Table_refsContext {
	var p = new(Table_refsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_table_refs

	return p
}

func (s *Table_refsContext) GetParser() antlr.Parser { return s.parser }

func (s *Table_refsContext) Name() antlr.TerminalNode {
	return s.GetToken(SyslParserName, 0)
}

func (s *Table_refsContext) TABLE() antlr.TerminalNode {
	return s.GetToken(SyslParserTABLE, 0)
}

func (s *Table_refsContext) TYPE() antlr.TerminalNode {
	return s.GetToken(SyslParserTYPE, 0)
}

func (s *Table_refsContext) Inplace_table_def() IInplace_table_defContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IInplace_table_defContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IInplace_table_defContext)
}

func (s *Table_refsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Table_refsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Table_refsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterTable_refs(s)
	}
}

func (s *Table_refsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitTable_refs(s)
	}
}

func (p *SyslParser) Table_refs() (localctx ITable_refsContext) {
	localctx = NewTable_refsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, SyslParserRULE_table_refs)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(336)
	_la = p.GetTokenStream().LA(1)

	if !(_la == SyslParserTABLE || _la == SyslParserTYPE) {
		p.GetErrorHandler().RecoverInline(p)
	} else {
		p.GetErrorHandler().ReportMatch(p)
		p.Consume()
	}
	{
		p.SetState(337)
		p.Match(SyslParserName)
	}
	p.SetState(339)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserCOLON {
		{
			p.SetState(338)
			p.Inplace_table_def()
		}

	}

	return localctx
}

// IFacadeContext is an interface to support dynamic dispatch.
type IFacadeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFacadeContext differentiates from other interfaces.
	IsFacadeContext()
}

type FacadeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFacadeContext() *FacadeContext {
	var p = new(FacadeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_facade
	return p
}

func (*FacadeContext) IsFacadeContext() {}

func NewFacadeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FacadeContext {
	var p = new(FacadeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_facade

	return p
}

func (s *FacadeContext) GetParser() antlr.Parser { return s.parser }

func (s *FacadeContext) WRAP() antlr.TerminalNode {
	return s.GetToken(SyslParserWRAP, 0)
}

func (s *FacadeContext) Model_name() IModel_nameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IModel_nameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IModel_nameContext)
}

func (s *FacadeContext) INDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserINDENT, 0)
}

func (s *FacadeContext) DEDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserDEDENT, 0)
}

func (s *FacadeContext) AllSYSL_COMMENT() []antlr.TerminalNode {
	return s.GetTokens(SyslParserSYSL_COMMENT)
}

func (s *FacadeContext) SYSL_COMMENT(i int) antlr.TerminalNode {
	return s.GetToken(SyslParserSYSL_COMMENT, i)
}

func (s *FacadeContext) AllTable_refs() []ITable_refsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ITable_refsContext)(nil)).Elem())
	var tst = make([]ITable_refsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ITable_refsContext)
		}
	}

	return tst
}

func (s *FacadeContext) Table_refs(i int) ITable_refsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITable_refsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ITable_refsContext)
}

func (s *FacadeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FacadeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FacadeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterFacade(s)
	}
}

func (s *FacadeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitFacade(s)
	}
}

func (p *SyslParser) Facade() (localctx IFacadeContext) {
	localctx = NewFacadeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, SyslParserRULE_facade)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(344)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SyslParserSYSL_COMMENT {
		{
			p.SetState(341)
			p.Match(SyslParserSYSL_COMMENT)
		}

		p.SetState(346)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(347)
		p.Match(SyslParserWRAP)
	}
	{
		p.SetState(348)
		p.Model_name()
	}
	{
		p.SetState(349)
		p.Match(SyslParserINDENT)
	}
	p.SetState(351)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == SyslParserTABLE || _la == SyslParserTYPE {
		{
			p.SetState(350)
			p.Table_refs()
		}

		p.SetState(353)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(355)
		p.Match(SyslParserDEDENT)
	}

	return localctx
}

// IDocumentation_stmtsContext is an interface to support dynamic dispatch.
type IDocumentation_stmtsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDocumentation_stmtsContext differentiates from other interfaces.
	IsDocumentation_stmtsContext()
}

type Documentation_stmtsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDocumentation_stmtsContext() *Documentation_stmtsContext {
	var p = new(Documentation_stmtsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_documentation_stmts
	return p
}

func (*Documentation_stmtsContext) IsDocumentation_stmtsContext() {}

func NewDocumentation_stmtsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Documentation_stmtsContext {
	var p = new(Documentation_stmtsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_documentation_stmts

	return p
}

func (s *Documentation_stmtsContext) GetParser() antlr.Parser { return s.parser }

func (s *Documentation_stmtsContext) AT() antlr.TerminalNode {
	return s.GetToken(SyslParserAT, 0)
}

func (s *Documentation_stmtsContext) Name() antlr.TerminalNode {
	return s.GetToken(SyslParserName, 0)
}

func (s *Documentation_stmtsContext) EQ() antlr.TerminalNode {
	return s.GetToken(SyslParserEQ, 0)
}

func (s *Documentation_stmtsContext) QSTRING() antlr.TerminalNode {
	return s.GetToken(SyslParserQSTRING, 0)
}

func (s *Documentation_stmtsContext) NEWLINE() antlr.TerminalNode {
	return s.GetToken(SyslParserNEWLINE, 0)
}

func (s *Documentation_stmtsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Documentation_stmtsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Documentation_stmtsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterDocumentation_stmts(s)
	}
}

func (s *Documentation_stmtsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitDocumentation_stmts(s)
	}
}

func (p *SyslParser) Documentation_stmts() (localctx IDocumentation_stmtsContext) {
	localctx = NewDocumentation_stmtsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, SyslParserRULE_documentation_stmts)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(357)
		p.Match(SyslParserAT)
	}
	{
		p.SetState(358)
		p.Match(SyslParserName)
	}
	{
		p.SetState(359)
		p.Match(SyslParserEQ)
	}
	{
		p.SetState(360)
		p.Match(SyslParserQSTRING)
	}
	{
		p.SetState(361)
		p.Match(SyslParserNEWLINE)
	}

	return localctx
}

// IVariable_substitutionContext is an interface to support dynamic dispatch.
type IVariable_substitutionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsVariable_substitutionContext differentiates from other interfaces.
	IsVariable_substitutionContext()
}

type Variable_substitutionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyVariable_substitutionContext() *Variable_substitutionContext {
	var p = new(Variable_substitutionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_variable_substitution
	return p
}

func (*Variable_substitutionContext) IsVariable_substitutionContext() {}

func NewVariable_substitutionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Variable_substitutionContext {
	var p = new(Variable_substitutionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_variable_substitution

	return p
}

func (s *Variable_substitutionContext) GetParser() antlr.Parser { return s.parser }

func (s *Variable_substitutionContext) FORWARD_SLASH() antlr.TerminalNode {
	return s.GetToken(SyslParserFORWARD_SLASH, 0)
}

func (s *Variable_substitutionContext) CURLY_OPEN() antlr.TerminalNode {
	return s.GetToken(SyslParserCURLY_OPEN, 0)
}

func (s *Variable_substitutionContext) CURLY_CLOSE() antlr.TerminalNode {
	return s.GetToken(SyslParserCURLY_CLOSE, 0)
}

func (s *Variable_substitutionContext) Field() IFieldContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFieldContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFieldContext)
}

func (s *Variable_substitutionContext) Name() antlr.TerminalNode {
	return s.GetToken(SyslParserName, 0)
}

func (s *Variable_substitutionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Variable_substitutionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Variable_substitutionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterVariable_substitution(s)
	}
}

func (s *Variable_substitutionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitVariable_substitution(s)
	}
}

func (p *SyslParser) Variable_substitution() (localctx IVariable_substitutionContext) {
	localctx = NewVariable_substitutionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, SyslParserRULE_variable_substitution)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(363)
		p.Match(SyslParserFORWARD_SLASH)
	}
	{
		p.SetState(364)
		p.Match(SyslParserCURLY_OPEN)
	}
	p.SetState(367)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 30, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(365)
			p.Field()
		}

	case 2:
		{
			p.SetState(366)
			p.Match(SyslParserName)
		}

	}
	{
		p.SetState(369)
		p.Match(SyslParserCURLY_CLOSE)
	}

	return localctx
}

// IStatic_pathContext is an interface to support dynamic dispatch.
type IStatic_pathContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStatic_pathContext differentiates from other interfaces.
	IsStatic_pathContext()
}

type Static_pathContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStatic_pathContext() *Static_pathContext {
	var p = new(Static_pathContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_static_path
	return p
}

func (*Static_pathContext) IsStatic_pathContext() {}

func NewStatic_pathContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Static_pathContext {
	var p = new(Static_pathContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_static_path

	return p
}

func (s *Static_pathContext) GetParser() antlr.Parser { return s.parser }

func (s *Static_pathContext) FORWARD_SLASH() antlr.TerminalNode {
	return s.GetToken(SyslParserFORWARD_SLASH, 0)
}

func (s *Static_pathContext) AllMINUS() []antlr.TerminalNode {
	return s.GetTokens(SyslParserMINUS)
}

func (s *Static_pathContext) MINUS(i int) antlr.TerminalNode {
	return s.GetToken(SyslParserMINUS, i)
}

func (s *Static_pathContext) AllName() []antlr.TerminalNode {
	return s.GetTokens(SyslParserName)
}

func (s *Static_pathContext) Name(i int) antlr.TerminalNode {
	return s.GetToken(SyslParserName, i)
}

func (s *Static_pathContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Static_pathContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Static_pathContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterStatic_path(s)
	}
}

func (s *Static_pathContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitStatic_path(s)
	}
}

func (p *SyslParser) Static_path() (localctx IStatic_pathContext) {
	localctx = NewStatic_pathContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 66, SyslParserRULE_static_path)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(371)
		p.Match(SyslParserFORWARD_SLASH)
	}
	p.SetState(373)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == SyslParserMINUS || _la == SyslParserName {
		p.SetState(372)
		_la = p.GetTokenStream().LA(1)

		if !(_la == SyslParserMINUS || _la == SyslParserName) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}

		p.SetState(375)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IQuery_varContext is an interface to support dynamic dispatch.
type IQuery_varContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsQuery_varContext differentiates from other interfaces.
	IsQuery_varContext()
}

type Query_varContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQuery_varContext() *Query_varContext {
	var p = new(Query_varContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_query_var
	return p
}

func (*Query_varContext) IsQuery_varContext() {}

func NewQuery_varContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Query_varContext {
	var p = new(Query_varContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_query_var

	return p
}

func (s *Query_varContext) GetParser() antlr.Parser { return s.parser }

func (s *Query_varContext) AllName() []antlr.TerminalNode {
	return s.GetTokens(SyslParserName)
}

func (s *Query_varContext) Name(i int) antlr.TerminalNode {
	return s.GetToken(SyslParserName, i)
}

func (s *Query_varContext) EQ() antlr.TerminalNode {
	return s.GetToken(SyslParserEQ, 0)
}

func (s *Query_varContext) NativeDataTypes() antlr.TerminalNode {
	return s.GetToken(SyslParserNativeDataTypes, 0)
}

func (s *Query_varContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Query_varContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Query_varContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterQuery_var(s)
	}
}

func (s *Query_varContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitQuery_var(s)
	}
}

func (p *SyslParser) Query_var() (localctx IQuery_varContext) {
	localctx = NewQuery_varContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 68, SyslParserRULE_query_var)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(377)
		p.Match(SyslParserName)
	}
	{
		p.SetState(378)
		p.Match(SyslParserEQ)
	}
	p.SetState(379)
	_la = p.GetTokenStream().LA(1)

	if !(_la == SyslParserNativeDataTypes || _la == SyslParserName) {
		p.GetErrorHandler().RecoverInline(p)
	} else {
		p.GetErrorHandler().ReportMatch(p)
		p.Consume()
	}

	return localctx
}

// IQuery_paramContext is an interface to support dynamic dispatch.
type IQuery_paramContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsQuery_paramContext differentiates from other interfaces.
	IsQuery_paramContext()
}

type Query_paramContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQuery_paramContext() *Query_paramContext {
	var p = new(Query_paramContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_query_param
	return p
}

func (*Query_paramContext) IsQuery_paramContext() {}

func NewQuery_paramContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Query_paramContext {
	var p = new(Query_paramContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_query_param

	return p
}

func (s *Query_paramContext) GetParser() antlr.Parser { return s.parser }

func (s *Query_paramContext) QN() antlr.TerminalNode {
	return s.GetToken(SyslParserQN, 0)
}

func (s *Query_paramContext) AllQuery_var() []IQuery_varContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IQuery_varContext)(nil)).Elem())
	var tst = make([]IQuery_varContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IQuery_varContext)
		}
	}

	return tst
}

func (s *Query_paramContext) Query_var(i int) IQuery_varContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IQuery_varContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IQuery_varContext)
}

func (s *Query_paramContext) AllAMP() []antlr.TerminalNode {
	return s.GetTokens(SyslParserAMP)
}

func (s *Query_paramContext) AMP(i int) antlr.TerminalNode {
	return s.GetToken(SyslParserAMP, i)
}

func (s *Query_paramContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Query_paramContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Query_paramContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterQuery_param(s)
	}
}

func (s *Query_paramContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitQuery_param(s)
	}
}

func (p *SyslParser) Query_param() (localctx IQuery_paramContext) {
	localctx = NewQuery_paramContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 70, SyslParserRULE_query_param)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(381)
		p.Match(SyslParserQN)
	}
	{
		p.SetState(382)
		p.Query_var()
	}
	p.SetState(387)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SyslParserAMP {
		{
			p.SetState(383)
			p.Match(SyslParserAMP)
		}
		{
			p.SetState(384)
			p.Query_var()
		}

		p.SetState(389)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IHttp_pathContext is an interface to support dynamic dispatch.
type IHttp_pathContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsHttp_pathContext differentiates from other interfaces.
	IsHttp_pathContext()
}

type Http_pathContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHttp_pathContext() *Http_pathContext {
	var p = new(Http_pathContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_http_path
	return p
}

func (*Http_pathContext) IsHttp_pathContext() {}

func NewHttp_pathContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Http_pathContext {
	var p = new(Http_pathContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_http_path

	return p
}

func (s *Http_pathContext) GetParser() antlr.Parser { return s.parser }

func (s *Http_pathContext) AllStatic_path() []IStatic_pathContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IStatic_pathContext)(nil)).Elem())
	var tst = make([]IStatic_pathContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IStatic_pathContext)
		}
	}

	return tst
}

func (s *Http_pathContext) Static_path(i int) IStatic_pathContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStatic_pathContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IStatic_pathContext)
}

func (s *Http_pathContext) AllVariable_substitution() []IVariable_substitutionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IVariable_substitutionContext)(nil)).Elem())
	var tst = make([]IVariable_substitutionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IVariable_substitutionContext)
		}
	}

	return tst
}

func (s *Http_pathContext) Variable_substitution(i int) IVariable_substitutionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVariable_substitutionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IVariable_substitutionContext)
}

func (s *Http_pathContext) Query_param() IQuery_paramContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IQuery_paramContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IQuery_paramContext)
}

func (s *Http_pathContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Http_pathContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Http_pathContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterHttp_path(s)
	}
}

func (s *Http_pathContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitHttp_path(s)
	}
}

func (p *SyslParser) Http_path() (localctx IHttp_pathContext) {
	localctx = NewHttp_pathContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 72, SyslParserRULE_http_path)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(392)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == SyslParserFORWARD_SLASH {
		p.SetState(392)
		p.GetErrorHandler().Sync(p)
		switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 33, p.GetParserRuleContext()) {
		case 1:
			{
				p.SetState(390)
				p.Static_path()
			}

		case 2:
			{
				p.SetState(391)
				p.Variable_substitution()
			}

		}

		p.SetState(394)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(397)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserQN {
		{
			p.SetState(396)
			p.Query_param()
		}

	}

	return localctx
}

// IEndpoint_nameContext is an interface to support dynamic dispatch.
type IEndpoint_nameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEndpoint_nameContext differentiates from other interfaces.
	IsEndpoint_nameContext()
}

type Endpoint_nameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEndpoint_nameContext() *Endpoint_nameContext {
	var p = new(Endpoint_nameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_endpoint_name
	return p
}

func (*Endpoint_nameContext) IsEndpoint_nameContext() {}

func NewEndpoint_nameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Endpoint_nameContext {
	var p = new(Endpoint_nameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_endpoint_name

	return p
}

func (s *Endpoint_nameContext) GetParser() antlr.Parser { return s.parser }

func (s *Endpoint_nameContext) COLON() antlr.TerminalNode {
	return s.GetToken(SyslParserCOLON, 0)
}

func (s *Endpoint_nameContext) Http_path() IHttp_pathContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHttp_pathContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IHttp_pathContext)
}

func (s *Endpoint_nameContext) TEXT_LINE() antlr.TerminalNode {
	return s.GetToken(SyslParserTEXT_LINE, 0)
}

func (s *Endpoint_nameContext) QSTRING() antlr.TerminalNode {
	return s.GetToken(SyslParserQSTRING, 0)
}

func (s *Endpoint_nameContext) Attribs_or_modifiers() IAttribs_or_modifiersContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAttribs_or_modifiersContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAttribs_or_modifiersContext)
}

func (s *Endpoint_nameContext) DOT() antlr.TerminalNode {
	return s.GetToken(SyslParserDOT, 0)
}

func (s *Endpoint_nameContext) AllName() []antlr.TerminalNode {
	return s.GetTokens(SyslParserName)
}

func (s *Endpoint_nameContext) Name(i int) antlr.TerminalNode {
	return s.GetToken(SyslParserName, i)
}

func (s *Endpoint_nameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Endpoint_nameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Endpoint_nameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterEndpoint_name(s)
	}
}

func (s *Endpoint_nameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitEndpoint_name(s)
	}
}

func (p *SyslParser) Endpoint_name() (localctx IEndpoint_nameContext) {
	localctx = NewEndpoint_nameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 74, SyslParserRULE_endpoint_name)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(409)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SyslParserFORWARD_SLASH:
		{
			p.SetState(399)
			p.Http_path()
		}

	case SyslParserDOT, SyslParserName:
		p.SetState(401)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SyslParserDOT {
			{
				p.SetState(400)
				p.Match(SyslParserDOT)
			}

		}
		p.SetState(404)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == SyslParserName {
			{
				p.SetState(403)
				p.Match(SyslParserName)
			}

			p.SetState(406)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case SyslParserTEXT_LINE:
		{
			p.SetState(408)
			p.Match(SyslParserTEXT_LINE)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	p.SetState(412)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserQSTRING {
		{
			p.SetState(411)
			p.Match(SyslParserQSTRING)
		}

	}
	p.SetState(415)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserSQ_OPEN {
		{
			p.SetState(414)
			p.Attribs_or_modifiers()
		}

	}
	{
		p.SetState(417)
		p.Match(SyslParserCOLON)
	}

	return localctx
}

// IRet_stmtContext is an interface to support dynamic dispatch.
type IRet_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsRet_stmtContext differentiates from other interfaces.
	IsRet_stmtContext()
}

type Ret_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRet_stmtContext() *Ret_stmtContext {
	var p = new(Ret_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_ret_stmt
	return p
}

func (*Ret_stmtContext) IsRet_stmtContext() {}

func NewRet_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Ret_stmtContext {
	var p = new(Ret_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_ret_stmt

	return p
}

func (s *Ret_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *Ret_stmtContext) RETURN() antlr.TerminalNode {
	return s.GetToken(SyslParserRETURN, 0)
}

func (s *Ret_stmtContext) TEXT() antlr.TerminalNode {
	return s.GetToken(SyslParserTEXT, 0)
}

func (s *Ret_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Ret_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Ret_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterRet_stmt(s)
	}
}

func (s *Ret_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitRet_stmt(s)
	}
}

func (p *SyslParser) Ret_stmt() (localctx IRet_stmtContext) {
	localctx = NewRet_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 76, SyslParserRULE_ret_stmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(419)
		p.Match(SyslParserRETURN)
	}
	{
		p.SetState(420)
		p.Match(SyslParserTEXT)
	}

	return localctx
}

// ITargetContext is an interface to support dynamic dispatch.
type ITargetContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTargetContext differentiates from other interfaces.
	IsTargetContext()
}

type TargetContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTargetContext() *TargetContext {
	var p = new(TargetContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_target
	return p
}

func (*TargetContext) IsTargetContext() {}

func NewTargetContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TargetContext {
	var p = new(TargetContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_target

	return p
}

func (s *TargetContext) GetParser() antlr.Parser { return s.parser }

func (s *TargetContext) DOT() antlr.TerminalNode {
	return s.GetToken(SyslParserDOT, 0)
}

func (s *TargetContext) App_name() IApp_nameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IApp_nameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IApp_nameContext)
}

func (s *TargetContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TargetContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TargetContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterTarget(s)
	}
}

func (s *TargetContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitTarget(s)
	}
}

func (p *SyslParser) Target() (localctx ITargetContext) {
	localctx = NewTargetContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 78, SyslParserRULE_target)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(424)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SyslParserDOT:
		{
			p.SetState(422)
			p.Match(SyslParserDOT)
		}

	case SyslParserTEXT_LINE, SyslParserName, SyslParserTEXT_NAME:
		{
			p.SetState(423)
			p.App_name()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// ITarget_endpointContext is an interface to support dynamic dispatch.
type ITarget_endpointContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTarget_endpointContext differentiates from other interfaces.
	IsTarget_endpointContext()
}

type Target_endpointContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTarget_endpointContext() *Target_endpointContext {
	var p = new(Target_endpointContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_target_endpoint
	return p
}

func (*Target_endpointContext) IsTarget_endpointContext() {}

func NewTarget_endpointContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Target_endpointContext {
	var p = new(Target_endpointContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_target_endpoint

	return p
}

func (s *Target_endpointContext) GetParser() antlr.Parser { return s.parser }

func (s *Target_endpointContext) TEXT_NAME() antlr.TerminalNode {
	return s.GetToken(SyslParserTEXT_NAME, 0)
}

func (s *Target_endpointContext) Name() antlr.TerminalNode {
	return s.GetToken(SyslParserName, 0)
}

func (s *Target_endpointContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Target_endpointContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Target_endpointContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterTarget_endpoint(s)
	}
}

func (s *Target_endpointContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitTarget_endpoint(s)
	}
}

func (p *SyslParser) Target_endpoint() (localctx ITarget_endpointContext) {
	localctx = NewTarget_endpointContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 80, SyslParserRULE_target_endpoint)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(426)
	_la = p.GetTokenStream().LA(1)

	if !(_la == SyslParserName || _la == SyslParserTEXT_NAME) {
		p.GetErrorHandler().RecoverInline(p)
	} else {
		p.GetErrorHandler().ReportMatch(p)
		p.Consume()
	}

	return localctx
}

// ICall_stmtContext is an interface to support dynamic dispatch.
type ICall_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCall_stmtContext differentiates from other interfaces.
	IsCall_stmtContext()
}

type Call_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCall_stmtContext() *Call_stmtContext {
	var p = new(Call_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_call_stmt
	return p
}

func (*Call_stmtContext) IsCall_stmtContext() {}

func NewCall_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Call_stmtContext {
	var p = new(Call_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_call_stmt

	return p
}

func (s *Call_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *Call_stmtContext) Target() ITargetContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITargetContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITargetContext)
}

func (s *Call_stmtContext) MEMBER() antlr.TerminalNode {
	return s.GetToken(SyslParserMEMBER, 0)
}

func (s *Call_stmtContext) Target_endpoint() ITarget_endpointContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITarget_endpointContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITarget_endpointContext)
}

func (s *Call_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Call_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Call_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterCall_stmt(s)
	}
}

func (s *Call_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitCall_stmt(s)
	}
}

func (p *SyslParser) Call_stmt() (localctx ICall_stmtContext) {
	localctx = NewCall_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 82, SyslParserRULE_call_stmt)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(428)
		p.Target()
	}
	p.SetState(431)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserMEMBER {
		{
			p.SetState(429)
			p.Match(SyslParserMEMBER)
		}
		{
			p.SetState(430)
			p.Target_endpoint()
		}

	}

	return localctx
}

// IIf_stmtContext is an interface to support dynamic dispatch.
type IIf_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsIf_stmtContext differentiates from other interfaces.
	IsIf_stmtContext()
}

type If_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIf_stmtContext() *If_stmtContext {
	var p = new(If_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_if_stmt
	return p
}

func (*If_stmtContext) IsIf_stmtContext() {}

func NewIf_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *If_stmtContext {
	var p = new(If_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_if_stmt

	return p
}

func (s *If_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *If_stmtContext) IF() antlr.TerminalNode {
	return s.GetToken(SyslParserIF, 0)
}

func (s *If_stmtContext) TEXT_NAME() antlr.TerminalNode {
	return s.GetToken(SyslParserTEXT_NAME, 0)
}

func (s *If_stmtContext) COLON() antlr.TerminalNode {
	return s.GetToken(SyslParserCOLON, 0)
}

func (s *If_stmtContext) INDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserINDENT, 0)
}

func (s *If_stmtContext) DEDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserDEDENT, 0)
}

func (s *If_stmtContext) AllHttp_statements() []IHttp_statementsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IHttp_statementsContext)(nil)).Elem())
	var tst = make([]IHttp_statementsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IHttp_statementsContext)
		}
	}

	return tst
}

func (s *If_stmtContext) Http_statements(i int) IHttp_statementsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHttp_statementsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IHttp_statementsContext)
}

func (s *If_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *If_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *If_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterIf_stmt(s)
	}
}

func (s *If_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitIf_stmt(s)
	}
}

func (p *SyslParser) If_stmt() (localctx IIf_stmtContext) {
	localctx = NewIf_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 84, SyslParserRULE_if_stmt)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(433)
		p.Match(SyslParserIF)
	}
	{
		p.SetState(434)
		p.Match(SyslParserTEXT_NAME)
	}
	{
		p.SetState(435)
		p.Match(SyslParserCOLON)
	}
	{
		p.SetState(436)
		p.Match(SyslParserINDENT)
	}
	p.SetState(440)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<SyslParserRETURN)|(1<<SyslParserIF)|(1<<SyslParserFOR)|(1<<SyslParserGROUP)|(1<<SyslParserWHATEVER)|(1<<SyslParserONE_OF))) != 0) || (((_la-33)&-(0x1f+1)) == 0 && ((1<<uint((_la-33)))&((1<<(SyslParserDOT-33))|(1<<(SyslParserAT-33))|(1<<(SyslParserPIPE-33))|(1<<(SyslParserQSTRING-33))|(1<<(SyslParserSYSL_COMMENT-33))|(1<<(SyslParserTEXT_LINE-33))|(1<<(SyslParserName-33))|(1<<(SyslParserTEXT_NAME-33)))) != 0) {
		{
			p.SetState(437)
			p.Http_statements()
		}

		p.SetState(442)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(443)
		p.Match(SyslParserDEDENT)
	}

	return localctx
}

// IIf_elseContext is an interface to support dynamic dispatch.
type IIf_elseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsIf_elseContext differentiates from other interfaces.
	IsIf_elseContext()
}

type If_elseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIf_elseContext() *If_elseContext {
	var p = new(If_elseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_if_else
	return p
}

func (*If_elseContext) IsIf_elseContext() {}

func NewIf_elseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *If_elseContext {
	var p = new(If_elseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_if_else

	return p
}

func (s *If_elseContext) GetParser() antlr.Parser { return s.parser }

func (s *If_elseContext) If_stmt() IIf_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIf_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIf_stmtContext)
}

func (s *If_elseContext) ELSE() antlr.TerminalNode {
	return s.GetToken(SyslParserELSE, 0)
}

func (s *If_elseContext) COLON() antlr.TerminalNode {
	return s.GetToken(SyslParserCOLON, 0)
}

func (s *If_elseContext) INDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserINDENT, 0)
}

func (s *If_elseContext) DEDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserDEDENT, 0)
}

func (s *If_elseContext) AllHttp_statements() []IHttp_statementsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IHttp_statementsContext)(nil)).Elem())
	var tst = make([]IHttp_statementsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IHttp_statementsContext)
		}
	}

	return tst
}

func (s *If_elseContext) Http_statements(i int) IHttp_statementsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHttp_statementsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IHttp_statementsContext)
}

func (s *If_elseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *If_elseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *If_elseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterIf_else(s)
	}
}

func (s *If_elseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitIf_else(s)
	}
}

func (p *SyslParser) If_else() (localctx IIf_elseContext) {
	localctx = NewIf_elseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 86, SyslParserRULE_if_else)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(445)
		p.If_stmt()
	}
	p.SetState(456)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserELSE {
		{
			p.SetState(446)
			p.Match(SyslParserELSE)
		}
		{
			p.SetState(447)
			p.Match(SyslParserCOLON)
		}
		{
			p.SetState(448)
			p.Match(SyslParserINDENT)
		}
		p.SetState(452)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<SyslParserRETURN)|(1<<SyslParserIF)|(1<<SyslParserFOR)|(1<<SyslParserGROUP)|(1<<SyslParserWHATEVER)|(1<<SyslParserONE_OF))) != 0) || (((_la-33)&-(0x1f+1)) == 0 && ((1<<uint((_la-33)))&((1<<(SyslParserDOT-33))|(1<<(SyslParserAT-33))|(1<<(SyslParserPIPE-33))|(1<<(SyslParserQSTRING-33))|(1<<(SyslParserSYSL_COMMENT-33))|(1<<(SyslParserTEXT_LINE-33))|(1<<(SyslParserName-33))|(1<<(SyslParserTEXT_NAME-33)))) != 0) {
			{
				p.SetState(449)
				p.Http_statements()
			}

			p.SetState(454)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(455)
			p.Match(SyslParserDEDENT)
		}

	}

	return localctx
}

// IFor_condContext is an interface to support dynamic dispatch.
type IFor_condContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFor_condContext differentiates from other interfaces.
	IsFor_condContext()
}

type For_condContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFor_condContext() *For_condContext {
	var p = new(For_condContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_for_cond
	return p
}

func (*For_condContext) IsFor_condContext() {}

func NewFor_condContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *For_condContext {
	var p = new(For_condContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_for_cond

	return p
}

func (s *For_condContext) GetParser() antlr.Parser { return s.parser }

func (s *For_condContext) TEXT_NAME() antlr.TerminalNode {
	return s.GetToken(SyslParserTEXT_NAME, 0)
}

func (s *For_condContext) COLON() antlr.TerminalNode {
	return s.GetToken(SyslParserCOLON, 0)
}

func (s *For_condContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *For_condContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *For_condContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterFor_cond(s)
	}
}

func (s *For_condContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitFor_cond(s)
	}
}

func (p *SyslParser) For_cond() (localctx IFor_condContext) {
	localctx = NewFor_condContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 88, SyslParserRULE_for_cond)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(458)
		p.Match(SyslParserTEXT_NAME)
	}
	{
		p.SetState(459)
		p.Match(SyslParserCOLON)
	}

	return localctx
}

// IFor_stmtContext is an interface to support dynamic dispatch.
type IFor_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFor_stmtContext differentiates from other interfaces.
	IsFor_stmtContext()
}

type For_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFor_stmtContext() *For_stmtContext {
	var p = new(For_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_for_stmt
	return p
}

func (*For_stmtContext) IsFor_stmtContext() {}

func NewFor_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *For_stmtContext {
	var p = new(For_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_for_stmt

	return p
}

func (s *For_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *For_stmtContext) FOR() antlr.TerminalNode {
	return s.GetToken(SyslParserFOR, 0)
}

func (s *For_stmtContext) For_cond() IFor_condContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFor_condContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFor_condContext)
}

func (s *For_stmtContext) INDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserINDENT, 0)
}

func (s *For_stmtContext) DEDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserDEDENT, 0)
}

func (s *For_stmtContext) AllHttp_statements() []IHttp_statementsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IHttp_statementsContext)(nil)).Elem())
	var tst = make([]IHttp_statementsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IHttp_statementsContext)
		}
	}

	return tst
}

func (s *For_stmtContext) Http_statements(i int) IHttp_statementsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHttp_statementsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IHttp_statementsContext)
}

func (s *For_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *For_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *For_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterFor_stmt(s)
	}
}

func (s *For_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitFor_stmt(s)
	}
}

func (p *SyslParser) For_stmt() (localctx IFor_stmtContext) {
	localctx = NewFor_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 90, SyslParserRULE_for_stmt)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(461)
		p.Match(SyslParserFOR)
	}
	{
		p.SetState(462)
		p.For_cond()
	}
	{
		p.SetState(463)
		p.Match(SyslParserINDENT)
	}
	p.SetState(467)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<SyslParserRETURN)|(1<<SyslParserIF)|(1<<SyslParserFOR)|(1<<SyslParserGROUP)|(1<<SyslParserWHATEVER)|(1<<SyslParserONE_OF))) != 0) || (((_la-33)&-(0x1f+1)) == 0 && ((1<<uint((_la-33)))&((1<<(SyslParserDOT-33))|(1<<(SyslParserAT-33))|(1<<(SyslParserPIPE-33))|(1<<(SyslParserQSTRING-33))|(1<<(SyslParserSYSL_COMMENT-33))|(1<<(SyslParserTEXT_LINE-33))|(1<<(SyslParserName-33))|(1<<(SyslParserTEXT_NAME-33)))) != 0) {
		{
			p.SetState(464)
			p.Http_statements()
		}

		p.SetState(469)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(470)
		p.Match(SyslParserDEDENT)
	}

	return localctx
}

// IHttp_method_commentContext is an interface to support dynamic dispatch.
type IHttp_method_commentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsHttp_method_commentContext differentiates from other interfaces.
	IsHttp_method_commentContext()
}

type Http_method_commentContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHttp_method_commentContext() *Http_method_commentContext {
	var p = new(Http_method_commentContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_http_method_comment
	return p
}

func (*Http_method_commentContext) IsHttp_method_commentContext() {}

func NewHttp_method_commentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Http_method_commentContext {
	var p = new(Http_method_commentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_http_method_comment

	return p
}

func (s *Http_method_commentContext) GetParser() antlr.Parser { return s.parser }

func (s *Http_method_commentContext) SYSL_COMMENT() antlr.TerminalNode {
	return s.GetToken(SyslParserSYSL_COMMENT, 0)
}

func (s *Http_method_commentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Http_method_commentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Http_method_commentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterHttp_method_comment(s)
	}
}

func (s *Http_method_commentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitHttp_method_comment(s)
	}
}

func (p *SyslParser) Http_method_comment() (localctx IHttp_method_commentContext) {
	localctx = NewHttp_method_commentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 92, SyslParserRULE_http_method_comment)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(472)
		p.Match(SyslParserSYSL_COMMENT)
	}

	return localctx
}

// IGroup_stmtContext is an interface to support dynamic dispatch.
type IGroup_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsGroup_stmtContext differentiates from other interfaces.
	IsGroup_stmtContext()
}

type Group_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGroup_stmtContext() *Group_stmtContext {
	var p = new(Group_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_group_stmt
	return p
}

func (*Group_stmtContext) IsGroup_stmtContext() {}

func NewGroup_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Group_stmtContext {
	var p = new(Group_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_group_stmt

	return p
}

func (s *Group_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *Group_stmtContext) GROUP() antlr.TerminalNode {
	return s.GetToken(SyslParserGROUP, 0)
}

func (s *Group_stmtContext) TEXT_NAME() antlr.TerminalNode {
	return s.GetToken(SyslParserTEXT_NAME, 0)
}

func (s *Group_stmtContext) COLON() antlr.TerminalNode {
	return s.GetToken(SyslParserCOLON, 0)
}

func (s *Group_stmtContext) INDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserINDENT, 0)
}

func (s *Group_stmtContext) DEDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserDEDENT, 0)
}

func (s *Group_stmtContext) AllHttp_statements() []IHttp_statementsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IHttp_statementsContext)(nil)).Elem())
	var tst = make([]IHttp_statementsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IHttp_statementsContext)
		}
	}

	return tst
}

func (s *Group_stmtContext) Http_statements(i int) IHttp_statementsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHttp_statementsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IHttp_statementsContext)
}

func (s *Group_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Group_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Group_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterGroup_stmt(s)
	}
}

func (s *Group_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitGroup_stmt(s)
	}
}

func (p *SyslParser) Group_stmt() (localctx IGroup_stmtContext) {
	localctx = NewGroup_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 94, SyslParserRULE_group_stmt)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(474)
		p.Match(SyslParserGROUP)
	}
	{
		p.SetState(475)
		p.Match(SyslParserTEXT_NAME)
	}
	{
		p.SetState(476)
		p.Match(SyslParserCOLON)
	}
	{
		p.SetState(477)
		p.Match(SyslParserINDENT)
	}
	p.SetState(479)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<SyslParserRETURN)|(1<<SyslParserIF)|(1<<SyslParserFOR)|(1<<SyslParserGROUP)|(1<<SyslParserWHATEVER)|(1<<SyslParserONE_OF))) != 0) || (((_la-33)&-(0x1f+1)) == 0 && ((1<<uint((_la-33)))&((1<<(SyslParserDOT-33))|(1<<(SyslParserAT-33))|(1<<(SyslParserPIPE-33))|(1<<(SyslParserQSTRING-33))|(1<<(SyslParserSYSL_COMMENT-33))|(1<<(SyslParserTEXT_LINE-33))|(1<<(SyslParserName-33))|(1<<(SyslParserTEXT_NAME-33)))) != 0) {
		{
			p.SetState(478)
			p.Http_statements()
		}

		p.SetState(481)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(483)
		p.Match(SyslParserDEDENT)
	}

	return localctx
}

// IOne_of_case_labelContext is an interface to support dynamic dispatch.
type IOne_of_case_labelContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsOne_of_case_labelContext differentiates from other interfaces.
	IsOne_of_case_labelContext()
}

type One_of_case_labelContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOne_of_case_labelContext() *One_of_case_labelContext {
	var p = new(One_of_case_labelContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_one_of_case_label
	return p
}

func (*One_of_case_labelContext) IsOne_of_case_labelContext() {}

func NewOne_of_case_labelContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *One_of_case_labelContext {
	var p = new(One_of_case_labelContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_one_of_case_label

	return p
}

func (s *One_of_case_labelContext) GetParser() antlr.Parser { return s.parser }

func (s *One_of_case_labelContext) AllName() []antlr.TerminalNode {
	return s.GetTokens(SyslParserName)
}

func (s *One_of_case_labelContext) Name(i int) antlr.TerminalNode {
	return s.GetToken(SyslParserName, i)
}

func (s *One_of_case_labelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *One_of_case_labelContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *One_of_case_labelContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterOne_of_case_label(s)
	}
}

func (s *One_of_case_labelContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitOne_of_case_label(s)
	}
}

func (p *SyslParser) One_of_case_label() (localctx IOne_of_case_labelContext) {
	localctx = NewOne_of_case_labelContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 96, SyslParserRULE_one_of_case_label)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(488)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SyslParserName {
		{
			p.SetState(485)
			p.Match(SyslParserName)
		}

		p.SetState(490)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IOne_of_casesContext is an interface to support dynamic dispatch.
type IOne_of_casesContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsOne_of_casesContext differentiates from other interfaces.
	IsOne_of_casesContext()
}

type One_of_casesContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOne_of_casesContext() *One_of_casesContext {
	var p = new(One_of_casesContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_one_of_cases
	return p
}

func (*One_of_casesContext) IsOne_of_casesContext() {}

func NewOne_of_casesContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *One_of_casesContext {
	var p = new(One_of_casesContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_one_of_cases

	return p
}

func (s *One_of_casesContext) GetParser() antlr.Parser { return s.parser }

func (s *One_of_casesContext) One_of_case_label() IOne_of_case_labelContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IOne_of_case_labelContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IOne_of_case_labelContext)
}

func (s *One_of_casesContext) COLON() antlr.TerminalNode {
	return s.GetToken(SyslParserCOLON, 0)
}

func (s *One_of_casesContext) INDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserINDENT, 0)
}

func (s *One_of_casesContext) DEDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserDEDENT, 0)
}

func (s *One_of_casesContext) AllHttp_statements() []IHttp_statementsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IHttp_statementsContext)(nil)).Elem())
	var tst = make([]IHttp_statementsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IHttp_statementsContext)
		}
	}

	return tst
}

func (s *One_of_casesContext) Http_statements(i int) IHttp_statementsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHttp_statementsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IHttp_statementsContext)
}

func (s *One_of_casesContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *One_of_casesContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *One_of_casesContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterOne_of_cases(s)
	}
}

func (s *One_of_casesContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitOne_of_cases(s)
	}
}

func (p *SyslParser) One_of_cases() (localctx IOne_of_casesContext) {
	localctx = NewOne_of_casesContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 98, SyslParserRULE_one_of_cases)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(491)
		p.One_of_case_label()
	}
	{
		p.SetState(492)
		p.Match(SyslParserCOLON)
	}
	{
		p.SetState(493)
		p.Match(SyslParserINDENT)
	}
	p.SetState(495)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<SyslParserRETURN)|(1<<SyslParserIF)|(1<<SyslParserFOR)|(1<<SyslParserGROUP)|(1<<SyslParserWHATEVER)|(1<<SyslParserONE_OF))) != 0) || (((_la-33)&-(0x1f+1)) == 0 && ((1<<uint((_la-33)))&((1<<(SyslParserDOT-33))|(1<<(SyslParserAT-33))|(1<<(SyslParserPIPE-33))|(1<<(SyslParserQSTRING-33))|(1<<(SyslParserSYSL_COMMENT-33))|(1<<(SyslParserTEXT_LINE-33))|(1<<(SyslParserName-33))|(1<<(SyslParserTEXT_NAME-33)))) != 0) {
		{
			p.SetState(494)
			p.Http_statements()
		}

		p.SetState(497)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(499)
		p.Match(SyslParserDEDENT)
	}

	return localctx
}

// IOne_of_stmtContext is an interface to support dynamic dispatch.
type IOne_of_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsOne_of_stmtContext differentiates from other interfaces.
	IsOne_of_stmtContext()
}

type One_of_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOne_of_stmtContext() *One_of_stmtContext {
	var p = new(One_of_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_one_of_stmt
	return p
}

func (*One_of_stmtContext) IsOne_of_stmtContext() {}

func NewOne_of_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *One_of_stmtContext {
	var p = new(One_of_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_one_of_stmt

	return p
}

func (s *One_of_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *One_of_stmtContext) ONE_OF() antlr.TerminalNode {
	return s.GetToken(SyslParserONE_OF, 0)
}

func (s *One_of_stmtContext) COLON() antlr.TerminalNode {
	return s.GetToken(SyslParserCOLON, 0)
}

func (s *One_of_stmtContext) INDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserINDENT, 0)
}

func (s *One_of_stmtContext) DEDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserDEDENT, 0)
}

func (s *One_of_stmtContext) AllOne_of_cases() []IOne_of_casesContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IOne_of_casesContext)(nil)).Elem())
	var tst = make([]IOne_of_casesContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IOne_of_casesContext)
		}
	}

	return tst
}

func (s *One_of_stmtContext) One_of_cases(i int) IOne_of_casesContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IOne_of_casesContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IOne_of_casesContext)
}

func (s *One_of_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *One_of_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *One_of_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterOne_of_stmt(s)
	}
}

func (s *One_of_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitOne_of_stmt(s)
	}
}

func (p *SyslParser) One_of_stmt() (localctx IOne_of_stmtContext) {
	localctx = NewOne_of_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 100, SyslParserRULE_one_of_stmt)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(501)
		p.Match(SyslParserONE_OF)
	}
	{
		p.SetState(502)
		p.Match(SyslParserCOLON)
	}
	{
		p.SetState(503)
		p.Match(SyslParserINDENT)
	}
	p.SetState(505)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == SyslParserCOLON || _la == SyslParserName {
		{
			p.SetState(504)
			p.One_of_cases()
		}

		p.SetState(507)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(509)
		p.Match(SyslParserDEDENT)
	}

	return localctx
}

// IText_stmtContext is an interface to support dynamic dispatch.
type IText_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsText_stmtContext differentiates from other interfaces.
	IsText_stmtContext()
}

type Text_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyText_stmtContext() *Text_stmtContext {
	var p = new(Text_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_text_stmt
	return p
}

func (*Text_stmtContext) IsText_stmtContext() {}

func NewText_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Text_stmtContext {
	var p = new(Text_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_text_stmt

	return p
}

func (s *Text_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *Text_stmtContext) TEXT_LINE() antlr.TerminalNode {
	return s.GetToken(SyslParserTEXT_LINE, 0)
}

func (s *Text_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Text_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Text_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterText_stmt(s)
	}
}

func (s *Text_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitText_stmt(s)
	}
}

func (p *SyslParser) Text_stmt() (localctx IText_stmtContext) {
	localctx = NewText_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 102, SyslParserRULE_text_stmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(511)
		p.Match(SyslParserTEXT_LINE)
	}

	return localctx
}

// IHttp_statementsContext is an interface to support dynamic dispatch.
type IHttp_statementsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsHttp_statementsContext differentiates from other interfaces.
	IsHttp_statementsContext()
}

type Http_statementsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHttp_statementsContext() *Http_statementsContext {
	var p = new(Http_statementsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_http_statements
	return p
}

func (*Http_statementsContext) IsHttp_statementsContext() {}

func NewHttp_statementsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Http_statementsContext {
	var p = new(Http_statementsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_http_statements

	return p
}

func (s *Http_statementsContext) GetParser() antlr.Parser { return s.parser }

func (s *Http_statementsContext) Doc_string() IDoc_stringContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDoc_stringContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDoc_stringContext)
}

func (s *Http_statementsContext) If_else() IIf_elseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIf_elseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIf_elseContext)
}

func (s *Http_statementsContext) For_stmt() IFor_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFor_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFor_stmtContext)
}

func (s *Http_statementsContext) Ret_stmt() IRet_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IRet_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IRet_stmtContext)
}

func (s *Http_statementsContext) Call_stmt() ICall_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICall_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICall_stmtContext)
}

func (s *Http_statementsContext) One_of_stmt() IOne_of_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IOne_of_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IOne_of_stmtContext)
}

func (s *Http_statementsContext) Http_method_comment() IHttp_method_commentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHttp_method_commentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IHttp_method_commentContext)
}

func (s *Http_statementsContext) Group_stmt() IGroup_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IGroup_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IGroup_stmtContext)
}

func (s *Http_statementsContext) QSTRING() antlr.TerminalNode {
	return s.GetToken(SyslParserQSTRING, 0)
}

func (s *Http_statementsContext) WHATEVER() antlr.TerminalNode {
	return s.GetToken(SyslParserWHATEVER, 0)
}

func (s *Http_statementsContext) Text_stmt() IText_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IText_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IText_stmtContext)
}

func (s *Http_statementsContext) Annotation() IAnnotationContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAnnotationContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAnnotationContext)
}

func (s *Http_statementsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Http_statementsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Http_statementsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterHttp_statements(s)
	}
}

func (s *Http_statementsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitHttp_statements(s)
	}
}

func (p *SyslParser) Http_statements() (localctx IHttp_statementsContext) {
	localctx = NewHttp_statementsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 104, SyslParserRULE_http_statements)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(525)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 51, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(513)
			p.Doc_string()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(514)
			p.If_else()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(515)
			p.For_stmt()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(516)
			p.Ret_stmt()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(517)
			p.Call_stmt()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(518)
			p.One_of_stmt()
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(519)
			p.Http_method_comment()
		}

	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(520)
			p.Group_stmt()
		}

	case 9:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(521)
			p.Match(SyslParserQSTRING)
		}

	case 10:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(522)
			p.Match(SyslParserWHATEVER)
		}

	case 11:
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(523)
			p.Text_stmt()
		}

	case 12:
		p.EnterOuterAlt(localctx, 12)
		{
			p.SetState(524)
			p.Annotation()
		}

	}

	return localctx
}

// IMethod_defContext is an interface to support dynamic dispatch.
type IMethod_defContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetUrl_path returns the url_path attribute.
	GetUrl_path() string

	// SetUrl_path sets the url_path attribute.
	SetUrl_path(string)

	// IsMethod_defContext differentiates from other interfaces.
	IsMethod_defContext()
}

type Method_defContext struct {
	*antlr.BaseParserRuleContext
	parser   antlr.Parser
	url_path string
}

func NewEmptyMethod_defContext() *Method_defContext {
	var p = new(Method_defContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_method_def
	return p
}

func (*Method_defContext) IsMethod_defContext() {}

func NewMethod_defContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int, url_path string) *Method_defContext {
	var p = new(Method_defContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_method_def

	p.url_path = url_path

	return p
}

func (s *Method_defContext) GetParser() antlr.Parser { return s.parser }

func (s *Method_defContext) GetUrl_path() string { return s.url_path }

func (s *Method_defContext) SetUrl_path(v string) { s.url_path = v }

func (s *Method_defContext) HTTP_VERBS() antlr.TerminalNode {
	return s.GetToken(SyslParserHTTP_VERBS, 0)
}

func (s *Method_defContext) COLON() antlr.TerminalNode {
	return s.GetToken(SyslParserCOLON, 0)
}

func (s *Method_defContext) INDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserINDENT, 0)
}

func (s *Method_defContext) DEDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserDEDENT, 0)
}

func (s *Method_defContext) Query_param() IQuery_paramContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IQuery_paramContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IQuery_paramContext)
}

func (s *Method_defContext) Attributes() IAttributesContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAttributesContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAttributesContext)
}

func (s *Method_defContext) AllHttp_statements() []IHttp_statementsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IHttp_statementsContext)(nil)).Elem())
	var tst = make([]IHttp_statementsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IHttp_statementsContext)
		}
	}

	return tst
}

func (s *Method_defContext) Http_statements(i int) IHttp_statementsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHttp_statementsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IHttp_statementsContext)
}

func (s *Method_defContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Method_defContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Method_defContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterMethod_def(s)
	}
}

func (s *Method_defContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitMethod_def(s)
	}
}

func (p *SyslParser) Method_def(url_path string) (localctx IMethod_defContext) {
	localctx = NewMethod_defContext(p, p.GetParserRuleContext(), p.GetState(), url_path)
	p.EnterRule(localctx, 106, SyslParserRULE_method_def)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(527)
		p.Match(SyslParserHTTP_VERBS)
	}
	p.SetState(529)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserQN {
		{
			p.SetState(528)
			p.Query_param()
		}

	}
	p.SetState(532)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserSQ_OPEN {
		{
			p.SetState(531)
			p.Attributes()
		}

	}
	{
		p.SetState(534)
		p.Match(SyslParserCOLON)
	}
	{
		p.SetState(535)
		p.Match(SyslParserINDENT)
	}
	p.SetState(539)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<SyslParserRETURN)|(1<<SyslParserIF)|(1<<SyslParserFOR)|(1<<SyslParserGROUP)|(1<<SyslParserWHATEVER)|(1<<SyslParserONE_OF))) != 0) || (((_la-33)&-(0x1f+1)) == 0 && ((1<<uint((_la-33)))&((1<<(SyslParserDOT-33))|(1<<(SyslParserAT-33))|(1<<(SyslParserPIPE-33))|(1<<(SyslParserQSTRING-33))|(1<<(SyslParserSYSL_COMMENT-33))|(1<<(SyslParserTEXT_LINE-33))|(1<<(SyslParserName-33))|(1<<(SyslParserTEXT_NAME-33)))) != 0) {
		{
			p.SetState(536)
			p.Http_statements()
		}

		p.SetState(541)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(542)
		p.Match(SyslParserDEDENT)
	}

	return localctx
}

// IEndpoint_declContext is an interface to support dynamic dispatch.
type IEndpoint_declContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetPrefix returns the prefix attribute.
	GetPrefix() string

	// SetPrefix sets the prefix attribute.
	SetPrefix(string)

	// IsEndpoint_declContext differentiates from other interfaces.
	IsEndpoint_declContext()
}

type Endpoint_declContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	prefix string
}

func NewEmptyEndpoint_declContext() *Endpoint_declContext {
	var p = new(Endpoint_declContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_endpoint_decl
	return p
}

func (*Endpoint_declContext) IsEndpoint_declContext() {}

func NewEndpoint_declContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int, prefix string) *Endpoint_declContext {
	var p = new(Endpoint_declContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_endpoint_decl

	p.prefix = prefix

	return p
}

func (s *Endpoint_declContext) GetParser() antlr.Parser { return s.parser }

func (s *Endpoint_declContext) GetPrefix() string { return s.prefix }

func (s *Endpoint_declContext) SetPrefix(v string) { s.prefix = v }

func (s *Endpoint_declContext) Api_endpoint() IApi_endpointContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IApi_endpointContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IApi_endpointContext)
}

func (s *Endpoint_declContext) Method_def() IMethod_defContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMethod_defContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMethod_defContext)
}

func (s *Endpoint_declContext) Http_statements() IHttp_statementsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHttp_statementsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IHttp_statementsContext)
}

func (s *Endpoint_declContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Endpoint_declContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Endpoint_declContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterEndpoint_decl(s)
	}
}

func (s *Endpoint_declContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitEndpoint_decl(s)
	}
}

func (p *SyslParser) Endpoint_decl(prefix string) (localctx IEndpoint_declContext) {
	localctx = NewEndpoint_declContext(p, p.GetParserRuleContext(), p.GetState(), prefix)
	p.EnterRule(localctx, 108, SyslParserRULE_endpoint_decl)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(547)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 55, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(544)
			p.Api_endpoint(localctx.(*Endpoint_declContext).prefix)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(545)
			p.Method_def(localctx.(*Endpoint_declContext).prefix)
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(546)
			p.Http_statements()
		}

	}

	return localctx
}

// IShortcutContext is an interface to support dynamic dispatch.
type IShortcutContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShortcutContext differentiates from other interfaces.
	IsShortcutContext()
}

type ShortcutContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShortcutContext() *ShortcutContext {
	var p = new(ShortcutContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_shortcut
	return p
}

func (*ShortcutContext) IsShortcutContext() {}

func NewShortcutContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShortcutContext {
	var p = new(ShortcutContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_shortcut

	return p
}

func (s *ShortcutContext) GetParser() antlr.Parser { return s.parser }

func (s *ShortcutContext) WHATEVER() antlr.TerminalNode {
	return s.GetToken(SyslParserWHATEVER, 0)
}

func (s *ShortcutContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShortcutContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ShortcutContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterShortcut(s)
	}
}

func (s *ShortcutContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitShortcut(s)
	}
}

func (p *SyslParser) Shortcut() (localctx IShortcutContext) {
	localctx = NewShortcutContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 110, SyslParserRULE_shortcut)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(549)
		p.Match(SyslParserWHATEVER)
	}

	return localctx
}

// IApi_endpointContext is an interface to support dynamic dispatch.
type IApi_endpointContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_endpoint_name returns the _endpoint_name rule contexts.
	Get_endpoint_name() IEndpoint_nameContext

	// Set_endpoint_name sets the _endpoint_name rule contexts.
	Set_endpoint_name(IEndpoint_nameContext)

	// GetPrefix returns the prefix attribute.
	GetPrefix() string

	// SetPrefix sets the prefix attribute.
	SetPrefix(string)

	// IsApi_endpointContext differentiates from other interfaces.
	IsApi_endpointContext()
}

type Api_endpointContext struct {
	*antlr.BaseParserRuleContext
	parser         antlr.Parser
	prefix         string
	_endpoint_name IEndpoint_nameContext
}

func NewEmptyApi_endpointContext() *Api_endpointContext {
	var p = new(Api_endpointContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_api_endpoint
	return p
}

func (*Api_endpointContext) IsApi_endpointContext() {}

func NewApi_endpointContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int, prefix string) *Api_endpointContext {
	var p = new(Api_endpointContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_api_endpoint

	p.prefix = prefix

	return p
}

func (s *Api_endpointContext) GetParser() antlr.Parser { return s.parser }

func (s *Api_endpointContext) Get_endpoint_name() IEndpoint_nameContext { return s._endpoint_name }

func (s *Api_endpointContext) Set_endpoint_name(v IEndpoint_nameContext) { s._endpoint_name = v }

func (s *Api_endpointContext) GetPrefix() string { return s.prefix }

func (s *Api_endpointContext) SetPrefix(v string) { s.prefix = v }

func (s *Api_endpointContext) WHATEVER() antlr.TerminalNode {
	return s.GetToken(SyslParserWHATEVER, 0)
}

func (s *Api_endpointContext) Endpoint_name() IEndpoint_nameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEndpoint_nameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IEndpoint_nameContext)
}

func (s *Api_endpointContext) Shortcut() IShortcutContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShortcutContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShortcutContext)
}

func (s *Api_endpointContext) INDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserINDENT, 0)
}

func (s *Api_endpointContext) DEDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserDEDENT, 0)
}

func (s *Api_endpointContext) AllEndpoint_decl() []IEndpoint_declContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IEndpoint_declContext)(nil)).Elem())
	var tst = make([]IEndpoint_declContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IEndpoint_declContext)
		}
	}

	return tst
}

func (s *Api_endpointContext) Endpoint_decl(i int) IEndpoint_declContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEndpoint_declContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IEndpoint_declContext)
}

func (s *Api_endpointContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Api_endpointContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Api_endpointContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterApi_endpoint(s)
	}
}

func (s *Api_endpointContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitApi_endpoint(s)
	}
}

func (p *SyslParser) Api_endpoint(prefix string) (localctx IApi_endpointContext) {
	localctx = NewApi_endpointContext(p, p.GetParserRuleContext(), p.GetState(), prefix)
	p.EnterRule(localctx, 112, SyslParserRULE_api_endpoint)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(564)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SyslParserWHATEVER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(551)
			p.Match(SyslParserWHATEVER)
		}

	case SyslParserFORWARD_SLASH, SyslParserDOT, SyslParserTEXT_LINE, SyslParserName:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(552)

			var _x = p.Endpoint_name()

			localctx.(*Api_endpointContext)._endpoint_name = _x
		}
		p.SetState(562)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case SyslParserWHATEVER:
			{
				p.SetState(553)
				p.Shortcut()
			}

		case SyslParserINDENT:
			{
				p.SetState(554)
				p.Match(SyslParserINDENT)
			}
			p.SetState(556)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			for ok := true; ok; ok = (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<SyslParserHTTP_VERBS)|(1<<SyslParserRETURN)|(1<<SyslParserIF)|(1<<SyslParserFOR)|(1<<SyslParserGROUP)|(1<<SyslParserWHATEVER)|(1<<SyslParserONE_OF)|(1<<SyslParserFORWARD_SLASH))) != 0) || (((_la-33)&-(0x1f+1)) == 0 && ((1<<uint((_la-33)))&((1<<(SyslParserDOT-33))|(1<<(SyslParserAT-33))|(1<<(SyslParserPIPE-33))|(1<<(SyslParserQSTRING-33))|(1<<(SyslParserSYSL_COMMENT-33))|(1<<(SyslParserTEXT_LINE-33))|(1<<(SyslParserName-33))|(1<<(SyslParserTEXT_NAME-33)))) != 0) {
				{
					p.SetState(555)
					p.Endpoint_decl(localctx.(*Api_endpointContext).prefix + (func() string {
						if localctx.(*Api_endpointContext).Get_endpoint_name() == nil {
							return ""
						} else {
							return p.GetTokenStream().GetTextFromTokens(localctx.(*Api_endpointContext).Get_endpoint_name().GetStart(), localctx.(*Api_endpointContext)._endpoint_name.GetStop())
						}
					}()))
				}

				p.SetState(558)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(560)
				p.Match(SyslParserDEDENT)
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// ICollector_stmtContext is an interface to support dynamic dispatch.
type ICollector_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCollector_stmtContext differentiates from other interfaces.
	IsCollector_stmtContext()
}

type Collector_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCollector_stmtContext() *Collector_stmtContext {
	var p = new(Collector_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_collector_stmt
	return p
}

func (*Collector_stmtContext) IsCollector_stmtContext() {}

func NewCollector_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Collector_stmtContext {
	var p = new(Collector_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_collector_stmt

	return p
}

func (s *Collector_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *Collector_stmtContext) Call_stmt() ICall_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICall_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICall_stmtContext)
}

func (s *Collector_stmtContext) HTTP_VERBS() antlr.TerminalNode {
	return s.GetToken(SyslParserHTTP_VERBS, 0)
}

func (s *Collector_stmtContext) Http_path() IHttp_pathContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHttp_pathContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IHttp_pathContext)
}

func (s *Collector_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Collector_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Collector_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterCollector_stmt(s)
	}
}

func (s *Collector_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitCollector_stmt(s)
	}
}

func (p *SyslParser) Collector_stmt() (localctx ICollector_stmtContext) {
	localctx = NewCollector_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 114, SyslParserRULE_collector_stmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(569)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SyslParserDOT, SyslParserTEXT_LINE, SyslParserName, SyslParserTEXT_NAME:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(566)
			p.Call_stmt()
		}

	case SyslParserHTTP_VERBS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(567)
			p.Match(SyslParserHTTP_VERBS)
		}
		{
			p.SetState(568)
			p.Http_path()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// ICollector_stmtsContext is an interface to support dynamic dispatch.
type ICollector_stmtsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCollector_stmtsContext differentiates from other interfaces.
	IsCollector_stmtsContext()
}

type Collector_stmtsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCollector_stmtsContext() *Collector_stmtsContext {
	var p = new(Collector_stmtsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_collector_stmts
	return p
}

func (*Collector_stmtsContext) IsCollector_stmtsContext() {}

func NewCollector_stmtsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Collector_stmtsContext {
	var p = new(Collector_stmtsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_collector_stmts

	return p
}

func (s *Collector_stmtsContext) GetParser() antlr.Parser { return s.parser }

func (s *Collector_stmtsContext) Collector_stmt() ICollector_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICollector_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICollector_stmtContext)
}

func (s *Collector_stmtsContext) Attribs_or_modifiers() IAttribs_or_modifiersContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAttribs_or_modifiersContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAttribs_or_modifiersContext)
}

func (s *Collector_stmtsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Collector_stmtsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Collector_stmtsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterCollector_stmts(s)
	}
}

func (s *Collector_stmtsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitCollector_stmts(s)
	}
}

func (p *SyslParser) Collector_stmts() (localctx ICollector_stmtsContext) {
	localctx = NewCollector_stmtsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 116, SyslParserRULE_collector_stmts)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(571)
		p.Collector_stmt()
	}
	{
		p.SetState(572)
		p.Attribs_or_modifiers()
	}

	return localctx
}

// ICollectorContext is an interface to support dynamic dispatch.
type ICollectorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCollectorContext differentiates from other interfaces.
	IsCollectorContext()
}

type CollectorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCollectorContext() *CollectorContext {
	var p = new(CollectorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_collector
	return p
}

func (*CollectorContext) IsCollectorContext() {}

func NewCollectorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CollectorContext {
	var p = new(CollectorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_collector

	return p
}

func (s *CollectorContext) GetParser() antlr.Parser { return s.parser }

func (s *CollectorContext) COLLECTOR() antlr.TerminalNode {
	return s.GetToken(SyslParserCOLLECTOR, 0)
}

func (s *CollectorContext) COLON() antlr.TerminalNode {
	return s.GetToken(SyslParserCOLON, 0)
}

func (s *CollectorContext) WHATEVER() antlr.TerminalNode {
	return s.GetToken(SyslParserWHATEVER, 0)
}

func (s *CollectorContext) INDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserINDENT, 0)
}

func (s *CollectorContext) DEDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserDEDENT, 0)
}

func (s *CollectorContext) AllCollector_stmts() []ICollector_stmtsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ICollector_stmtsContext)(nil)).Elem())
	var tst = make([]ICollector_stmtsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ICollector_stmtsContext)
		}
	}

	return tst
}

func (s *CollectorContext) Collector_stmts(i int) ICollector_stmtsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICollector_stmtsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ICollector_stmtsContext)
}

func (s *CollectorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CollectorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CollectorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterCollector(s)
	}
}

func (s *CollectorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitCollector(s)
	}
}

func (p *SyslParser) Collector() (localctx ICollectorContext) {
	localctx = NewCollectorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 118, SyslParserRULE_collector)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(574)
		p.Match(SyslParserCOLLECTOR)
	}
	{
		p.SetState(575)
		p.Match(SyslParserCOLON)
	}
	p.SetState(585)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SyslParserWHATEVER:
		{
			p.SetState(576)
			p.Match(SyslParserWHATEVER)
		}

	case SyslParserINDENT:
		{
			p.SetState(577)
			p.Match(SyslParserINDENT)
		}
		p.SetState(579)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == SyslParserHTTP_VERBS || (((_la-33)&-(0x1f+1)) == 0 && ((1<<uint((_la-33)))&((1<<(SyslParserDOT-33))|(1<<(SyslParserTEXT_LINE-33))|(1<<(SyslParserName-33))|(1<<(SyslParserTEXT_NAME-33)))) != 0) {
			{
				p.SetState(578)
				p.Collector_stmts()
			}

			p.SetState(581)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(583)
			p.Match(SyslParserDEDENT)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEventContext is an interface to support dynamic dispatch.
type IEventContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEventContext differentiates from other interfaces.
	IsEventContext()
}

type EventContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventContext() *EventContext {
	var p = new(EventContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_event
	return p
}

func (*EventContext) IsEventContext() {}

func NewEventContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventContext {
	var p = new(EventContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_event

	return p
}

func (s *EventContext) GetParser() antlr.Parser { return s.parser }

func (s *EventContext) DISTANCE() antlr.TerminalNode {
	return s.GetToken(SyslParserDISTANCE, 0)
}

func (s *EventContext) EVENT_NAME() antlr.TerminalNode {
	return s.GetToken(SyslParserEVENT_NAME, 0)
}

func (s *EventContext) COLON() antlr.TerminalNode {
	return s.GetToken(SyslParserCOLON, 0)
}

func (s *EventContext) WHATEVER() antlr.TerminalNode {
	return s.GetToken(SyslParserWHATEVER, 0)
}

func (s *EventContext) INDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserINDENT, 0)
}

func (s *EventContext) DEDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserDEDENT, 0)
}

func (s *EventContext) Attribs_or_modifiers() IAttribs_or_modifiersContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAttribs_or_modifiersContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAttribs_or_modifiersContext)
}

func (s *EventContext) AllHttp_statements() []IHttp_statementsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IHttp_statementsContext)(nil)).Elem())
	var tst = make([]IHttp_statementsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IHttp_statementsContext)
		}
	}

	return tst
}

func (s *EventContext) Http_statements(i int) IHttp_statementsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHttp_statementsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IHttp_statementsContext)
}

func (s *EventContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterEvent(s)
	}
}

func (s *EventContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitEvent(s)
	}
}

func (p *SyslParser) Event() (localctx IEventContext) {
	localctx = NewEventContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 120, SyslParserRULE_event)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(587)
		p.Match(SyslParserDISTANCE)
	}
	{
		p.SetState(588)
		p.Match(SyslParserEVENT_NAME)
	}
	p.SetState(590)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserSQ_OPEN {
		{
			p.SetState(589)
			p.Attribs_or_modifiers()
		}

	}
	{
		p.SetState(592)
		p.Match(SyslParserCOLON)
	}
	p.SetState(602)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SyslParserWHATEVER:
		{
			p.SetState(593)
			p.Match(SyslParserWHATEVER)
		}

	case SyslParserINDENT:
		{
			p.SetState(594)
			p.Match(SyslParserINDENT)
		}
		p.SetState(596)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<SyslParserRETURN)|(1<<SyslParserIF)|(1<<SyslParserFOR)|(1<<SyslParserGROUP)|(1<<SyslParserWHATEVER)|(1<<SyslParserONE_OF))) != 0) || (((_la-33)&-(0x1f+1)) == 0 && ((1<<uint((_la-33)))&((1<<(SyslParserDOT-33))|(1<<(SyslParserAT-33))|(1<<(SyslParserPIPE-33))|(1<<(SyslParserQSTRING-33))|(1<<(SyslParserSYSL_COMMENT-33))|(1<<(SyslParserTEXT_LINE-33))|(1<<(SyslParserName-33))|(1<<(SyslParserTEXT_NAME-33)))) != 0) {
			{
				p.SetState(595)
				p.Http_statements()
			}

			p.SetState(598)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(600)
			p.Match(SyslParserDEDENT)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IApp_declContext is an interface to support dynamic dispatch.
type IApp_declContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsApp_declContext differentiates from other interfaces.
	IsApp_declContext()
}

type App_declContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyApp_declContext() *App_declContext {
	var p = new(App_declContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_app_decl
	return p
}

func (*App_declContext) IsApp_declContext() {}

func NewApp_declContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *App_declContext {
	var p = new(App_declContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_app_decl

	return p
}

func (s *App_declContext) GetParser() antlr.Parser { return s.parser }

func (s *App_declContext) INDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserINDENT, 0)
}

func (s *App_declContext) DEDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserDEDENT, 0)
}

func (s *App_declContext) AllTable() []ITableContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ITableContext)(nil)).Elem())
	var tst = make([]ITableContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ITableContext)
		}
	}

	return tst
}

func (s *App_declContext) Table(i int) ITableContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITableContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ITableContext)
}

func (s *App_declContext) AllFacade() []IFacadeContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IFacadeContext)(nil)).Elem())
	var tst = make([]IFacadeContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IFacadeContext)
		}
	}

	return tst
}

func (s *App_declContext) Facade(i int) IFacadeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFacadeContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IFacadeContext)
}

func (s *App_declContext) AllSYSL_COMMENT() []antlr.TerminalNode {
	return s.GetTokens(SyslParserSYSL_COMMENT)
}

func (s *App_declContext) SYSL_COMMENT(i int) antlr.TerminalNode {
	return s.GetToken(SyslParserSYSL_COMMENT, i)
}

func (s *App_declContext) AllApi_endpoint() []IApi_endpointContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IApi_endpointContext)(nil)).Elem())
	var tst = make([]IApi_endpointContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IApi_endpointContext)
		}
	}

	return tst
}

func (s *App_declContext) Api_endpoint(i int) IApi_endpointContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IApi_endpointContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IApi_endpointContext)
}

func (s *App_declContext) AllCollector() []ICollectorContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ICollectorContext)(nil)).Elem())
	var tst = make([]ICollectorContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ICollectorContext)
		}
	}

	return tst
}

func (s *App_declContext) Collector(i int) ICollectorContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICollectorContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ICollectorContext)
}

func (s *App_declContext) AllEvent() []IEventContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IEventContext)(nil)).Elem())
	var tst = make([]IEventContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IEventContext)
		}
	}

	return tst
}

func (s *App_declContext) Event(i int) IEventContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEventContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IEventContext)
}

func (s *App_declContext) AllAnnotation() []IAnnotationContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IAnnotationContext)(nil)).Elem())
	var tst = make([]IAnnotationContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IAnnotationContext)
		}
	}

	return tst
}

func (s *App_declContext) Annotation(i int) IAnnotationContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAnnotationContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IAnnotationContext)
}

func (s *App_declContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *App_declContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *App_declContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterApp_decl(s)
	}
}

func (s *App_declContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitApp_decl(s)
	}
}

func (p *SyslParser) App_decl() (localctx IApp_declContext) {
	localctx = NewApp_declContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 122, SyslParserRULE_app_decl)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(604)
		p.Match(SyslParserINDENT)
	}
	p.SetState(612)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<SyslParserWRAP)|(1<<SyslParserTABLE)|(1<<SyslParserTYPE)|(1<<SyslParserWHATEVER)|(1<<SyslParserDISTANCE)|(1<<SyslParserCOLLECTOR)|(1<<SyslParserFORWARD_SLASH))) != 0) || (((_la-33)&-(0x1f+1)) == 0 && ((1<<uint((_la-33)))&((1<<(SyslParserDOT-33))|(1<<(SyslParserAT-33))|(1<<(SyslParserSYSL_COMMENT-33))|(1<<(SyslParserTEXT_LINE-33))|(1<<(SyslParserName-33)))) != 0) {
		p.SetState(612)
		p.GetErrorHandler().Sync(p)
		switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 65, p.GetParserRuleContext()) {
		case 1:
			{
				p.SetState(605)
				p.Table()
			}

		case 2:
			{
				p.SetState(606)
				p.Facade()
			}

		case 3:
			{
				p.SetState(607)
				p.Match(SyslParserSYSL_COMMENT)
			}

		case 4:
			{
				p.SetState(608)
				p.Api_endpoint("")
			}

		case 5:
			{
				p.SetState(609)
				p.Collector()
			}

		case 6:
			{
				p.SetState(610)
				p.Event()
			}

		case 7:
			{
				p.SetState(611)
				p.Annotation()
			}

		}

		p.SetState(614)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(616)
		p.Match(SyslParserDEDENT)
	}

	return localctx
}

// IApplicationContext is an interface to support dynamic dispatch.
type IApplicationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsApplicationContext differentiates from other interfaces.
	IsApplicationContext()
}

type ApplicationContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyApplicationContext() *ApplicationContext {
	var p = new(ApplicationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_application
	return p
}

func (*ApplicationContext) IsApplicationContext() {}

func NewApplicationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ApplicationContext {
	var p = new(ApplicationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_application

	return p
}

func (s *ApplicationContext) GetParser() antlr.Parser { return s.parser }

func (s *ApplicationContext) Name_with_attribs() IName_with_attribsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IName_with_attribsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IName_with_attribsContext)
}

func (s *ApplicationContext) COLON() antlr.TerminalNode {
	return s.GetToken(SyslParserCOLON, 0)
}

func (s *ApplicationContext) App_decl() IApp_declContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IApp_declContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IApp_declContext)
}

func (s *ApplicationContext) AllSYSL_COMMENT() []antlr.TerminalNode {
	return s.GetTokens(SyslParserSYSL_COMMENT)
}

func (s *ApplicationContext) SYSL_COMMENT(i int) antlr.TerminalNode {
	return s.GetToken(SyslParserSYSL_COMMENT, i)
}

func (s *ApplicationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ApplicationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ApplicationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterApplication(s)
	}
}

func (s *ApplicationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitApplication(s)
	}
}

func (p *SyslParser) Application() (localctx IApplicationContext) {
	localctx = NewApplicationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 124, SyslParserRULE_application)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(621)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SyslParserSYSL_COMMENT {
		{
			p.SetState(618)
			p.Match(SyslParserSYSL_COMMENT)
		}

		p.SetState(623)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(624)
		p.Name_with_attribs()
	}
	{
		p.SetState(625)
		p.Match(SyslParserCOLON)
	}
	{
		p.SetState(626)
		p.App_decl()
	}

	return localctx
}

// IPathContext is an interface to support dynamic dispatch.
type IPathContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPathContext differentiates from other interfaces.
	IsPathContext()
}

type PathContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPathContext() *PathContext {
	var p = new(PathContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_path
	return p
}

func (*PathContext) IsPathContext() {}

func NewPathContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PathContext {
	var p = new(PathContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_path

	return p
}

func (s *PathContext) GetParser() antlr.Parser { return s.parser }

func (s *PathContext) AllName() []antlr.TerminalNode {
	return s.GetTokens(SyslParserName)
}

func (s *PathContext) Name(i int) antlr.TerminalNode {
	return s.GetToken(SyslParserName, i)
}

func (s *PathContext) AllFORWARD_SLASH() []antlr.TerminalNode {
	return s.GetTokens(SyslParserFORWARD_SLASH)
}

func (s *PathContext) FORWARD_SLASH(i int) antlr.TerminalNode {
	return s.GetToken(SyslParserFORWARD_SLASH, i)
}

func (s *PathContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PathContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PathContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterPath(s)
	}
}

func (s *PathContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitPath(s)
	}
}

func (p *SyslParser) Path() (localctx IPathContext) {
	localctx = NewPathContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 126, SyslParserRULE_path)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(629)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserFORWARD_SLASH {
		{
			p.SetState(628)
			p.Match(SyslParserFORWARD_SLASH)
		}

	}
	{
		p.SetState(631)
		p.Match(SyslParserName)
	}
	p.SetState(636)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SyslParserFORWARD_SLASH {
		{
			p.SetState(632)
			p.Match(SyslParserFORWARD_SLASH)
		}
		{
			p.SetState(633)
			p.Match(SyslParserName)
		}

		p.SetState(638)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IImport_stmtContext is an interface to support dynamic dispatch.
type IImport_stmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsImport_stmtContext differentiates from other interfaces.
	IsImport_stmtContext()
}

type Import_stmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyImport_stmtContext() *Import_stmtContext {
	var p = new(Import_stmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_import_stmt
	return p
}

func (*Import_stmtContext) IsImport_stmtContext() {}

func NewImport_stmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Import_stmtContext {
	var p = new(Import_stmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_import_stmt

	return p
}

func (s *Import_stmtContext) GetParser() antlr.Parser { return s.parser }

func (s *Import_stmtContext) IMPORT() antlr.TerminalNode {
	return s.GetToken(SyslParserIMPORT, 0)
}

func (s *Import_stmtContext) TEXT() antlr.TerminalNode {
	return s.GetToken(SyslParserTEXT, 0)
}

func (s *Import_stmtContext) AllSYSL_COMMENT() []antlr.TerminalNode {
	return s.GetTokens(SyslParserSYSL_COMMENT)
}

func (s *Import_stmtContext) SYSL_COMMENT(i int) antlr.TerminalNode {
	return s.GetToken(SyslParserSYSL_COMMENT, i)
}

func (s *Import_stmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Import_stmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Import_stmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterImport_stmt(s)
	}
}

func (s *Import_stmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitImport_stmt(s)
	}
}

func (p *SyslParser) Import_stmt() (localctx IImport_stmtContext) {
	localctx = NewImport_stmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 128, SyslParserRULE_import_stmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(639)
		p.Match(SyslParserIMPORT)
	}
	{
		p.SetState(640)
		p.Match(SyslParserTEXT)
	}
	p.SetState(644)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 70, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(641)
				p.Match(SyslParserSYSL_COMMENT)
			}

		}
		p.SetState(646)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 70, p.GetParserRuleContext())
	}

	return localctx
}

// IImports_declContext is an interface to support dynamic dispatch.
type IImports_declContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsImports_declContext differentiates from other interfaces.
	IsImports_declContext()
}

type Imports_declContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyImports_declContext() *Imports_declContext {
	var p = new(Imports_declContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_imports_decl
	return p
}

func (*Imports_declContext) IsImports_declContext() {}

func NewImports_declContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Imports_declContext {
	var p = new(Imports_declContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_imports_decl

	return p
}

func (s *Imports_declContext) GetParser() antlr.Parser { return s.parser }

func (s *Imports_declContext) AllImport_stmt() []IImport_stmtContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IImport_stmtContext)(nil)).Elem())
	var tst = make([]IImport_stmtContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IImport_stmtContext)
		}
	}

	return tst
}

func (s *Imports_declContext) Import_stmt(i int) IImport_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IImport_stmtContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IImport_stmtContext)
}

func (s *Imports_declContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Imports_declContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Imports_declContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterImports_decl(s)
	}
}

func (s *Imports_declContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitImports_decl(s)
	}
}

func (p *SyslParser) Imports_decl() (localctx IImports_declContext) {
	localctx = NewImports_declContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 130, SyslParserRULE_imports_decl)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(648)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == SyslParserIMPORT {
		{
			p.SetState(647)
			p.Import_stmt()
		}

		p.SetState(650)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// ISysl_fileContext is an interface to support dynamic dispatch.
type ISysl_fileContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSysl_fileContext differentiates from other interfaces.
	IsSysl_fileContext()
}

type Sysl_fileContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySysl_fileContext() *Sysl_fileContext {
	var p = new(Sysl_fileContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_sysl_file
	return p
}

func (*Sysl_fileContext) IsSysl_fileContext() {}

func NewSysl_fileContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Sysl_fileContext {
	var p = new(Sysl_fileContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_sysl_file

	return p
}

func (s *Sysl_fileContext) GetParser() antlr.Parser { return s.parser }

func (s *Sysl_fileContext) EOF() antlr.TerminalNode {
	return s.GetToken(SyslParserEOF, 0)
}

func (s *Sysl_fileContext) Imports_decl() IImports_declContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IImports_declContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IImports_declContext)
}

func (s *Sysl_fileContext) AllApplication() []IApplicationContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IApplicationContext)(nil)).Elem())
	var tst = make([]IApplicationContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IApplicationContext)
		}
	}

	return tst
}

func (s *Sysl_fileContext) Application(i int) IApplicationContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IApplicationContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IApplicationContext)
}

func (s *Sysl_fileContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Sysl_fileContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Sysl_fileContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterSysl_file(s)
	}
}

func (s *Sysl_fileContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitSysl_file(s)
	}
}

func (p *SyslParser) Sysl_file() (localctx ISysl_fileContext) {
	localctx = NewSysl_fileContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 132, SyslParserRULE_sysl_file)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(653)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserIMPORT {
		{
			p.SetState(652)
			p.Imports_decl()
		}

	}
	p.SetState(656)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (((_la-55)&-(0x1f+1)) == 0 && ((1<<uint((_la-55)))&((1<<(SyslParserSYSL_COMMENT-55))|(1<<(SyslParserTEXT_LINE-55))|(1<<(SyslParserName-55))|(1<<(SyslParserTEXT_NAME-55)))) != 0) {
		{
			p.SetState(655)
			p.Application()
		}

		p.SetState(658)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(660)
		p.Match(SyslParserEOF)
	}

	return localctx
}
