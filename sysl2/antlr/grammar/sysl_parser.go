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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 65, 696,
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
	4, 66, 9, 66, 4, 67, 9, 67, 4, 68, 9, 68, 4, 69, 9, 69, 4, 70, 9, 70, 3,
	2, 3, 2, 3, 2, 3, 2, 7, 2, 145, 10, 2, 12, 2, 14, 2, 148, 11, 2, 3, 3,
	3, 3, 3, 3, 3, 3, 5, 3, 154, 10, 3, 3, 3, 3, 3, 3, 4, 3, 4, 3, 4, 7, 4,
	161, 10, 4, 12, 4, 14, 4, 164, 11, 4, 3, 5, 3, 5, 3, 5, 3, 5, 3, 6, 3,
	6, 3, 6, 3, 6, 3, 7, 3, 7, 3, 7, 3, 8, 3, 8, 3, 9, 3, 9, 3, 9, 3, 9, 7,
	9, 183, 10, 9, 12, 9, 14, 9, 186, 11, 9, 3, 9, 3, 9, 3, 10, 3, 10, 3, 10,
	3, 10, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 5, 11, 199, 10, 11, 3, 12, 3,
	12, 3, 12, 3, 12, 7, 12, 205, 10, 12, 12, 12, 14, 12, 208, 11, 12, 3, 12,
	3, 12, 3, 13, 3, 13, 5, 13, 214, 10, 13, 3, 14, 3, 14, 3, 14, 3, 14, 7,
	14, 220, 10, 14, 12, 14, 14, 14, 223, 11, 14, 3, 14, 3, 14, 3, 15, 3, 15,
	3, 15, 5, 15, 230, 10, 15, 3, 16, 3, 16, 3, 17, 3, 17, 3, 18, 3, 18, 3,
	18, 6, 18, 239, 10, 18, 13, 18, 14, 18, 240, 3, 18, 3, 18, 3, 19, 3, 19,
	3, 19, 5, 19, 248, 10, 19, 3, 20, 3, 20, 3, 20, 3, 20, 3, 20, 3, 21, 3,
	21, 6, 21, 257, 10, 21, 13, 21, 14, 21, 258, 3, 21, 3, 21, 3, 22, 3, 22,
	3, 22, 3, 22, 5, 22, 267, 10, 22, 3, 22, 5, 22, 270, 10, 22, 3, 22, 5,
	22, 273, 10, 22, 3, 22, 5, 22, 276, 10, 22, 3, 22, 3, 22, 5, 22, 280, 10,
	22, 5, 22, 282, 10, 22, 3, 23, 3, 23, 3, 23, 3, 23, 3, 24, 7, 24, 289,
	10, 24, 12, 24, 14, 24, 292, 11, 24, 3, 24, 3, 24, 3, 24, 5, 24, 297, 10,
	24, 3, 24, 3, 24, 3, 24, 3, 24, 6, 24, 303, 10, 24, 13, 24, 14, 24, 304,
	3, 24, 3, 24, 3, 25, 3, 25, 3, 25, 5, 25, 312, 10, 25, 3, 25, 3, 25, 5,
	25, 316, 10, 25, 3, 26, 3, 26, 3, 26, 3, 27, 3, 27, 7, 27, 323, 10, 27,
	12, 27, 14, 27, 326, 11, 27, 3, 28, 3, 28, 5, 28, 330, 10, 28, 3, 28, 5,
	28, 333, 10, 28, 3, 29, 3, 29, 3, 29, 3, 30, 3, 30, 3, 30, 3, 30, 5, 30,
	342, 10, 30, 6, 30, 344, 10, 30, 13, 30, 14, 30, 345, 3, 30, 3, 30, 3,
	31, 3, 31, 3, 31, 5, 31, 353, 10, 31, 3, 32, 7, 32, 356, 10, 32, 12, 32,
	14, 32, 359, 11, 32, 3, 32, 3, 32, 3, 32, 3, 32, 6, 32, 365, 10, 32, 13,
	32, 14, 32, 366, 3, 32, 3, 32, 3, 33, 3, 33, 3, 33, 3, 33, 3, 33, 3, 33,
	3, 34, 3, 34, 3, 34, 3, 34, 3, 35, 3, 35, 3, 35, 3, 35, 3, 35, 5, 35, 386,
	10, 35, 3, 36, 3, 36, 3, 36, 3, 36, 7, 36, 392, 10, 36, 12, 36, 14, 36,
	395, 11, 36, 3, 37, 3, 37, 3, 37, 3, 37, 3, 37, 3, 37, 3, 38, 3, 38, 3,
	39, 3, 39, 3, 39, 5, 39, 408, 10, 39, 3, 40, 3, 40, 6, 40, 412, 10, 40,
	13, 40, 14, 40, 413, 5, 40, 416, 10, 40, 3, 41, 5, 41, 419, 10, 41, 3,
	41, 6, 41, 422, 10, 41, 13, 41, 14, 41, 423, 3, 41, 5, 41, 427, 10, 41,
	3, 42, 3, 42, 3, 42, 3, 43, 3, 43, 5, 43, 434, 10, 43, 3, 44, 3, 44, 3,
	45, 3, 45, 3, 45, 3, 45, 3, 46, 3, 46, 3, 46, 3, 46, 3, 46, 7, 46, 447,
	10, 46, 12, 46, 14, 46, 450, 11, 46, 3, 46, 3, 46, 3, 47, 3, 47, 3, 47,
	3, 47, 3, 47, 7, 47, 459, 10, 47, 12, 47, 14, 47, 462, 11, 47, 3, 47, 5,
	47, 465, 10, 47, 3, 48, 3, 48, 3, 48, 3, 49, 3, 49, 3, 49, 3, 49, 7, 49,
	474, 10, 49, 12, 49, 14, 49, 477, 11, 49, 3, 49, 3, 49, 3, 50, 3, 50, 3,
	51, 7, 51, 484, 10, 51, 12, 51, 14, 51, 487, 11, 51, 3, 52, 3, 52, 3, 52,
	3, 52, 6, 52, 493, 10, 52, 13, 52, 14, 52, 494, 3, 52, 3, 52, 3, 53, 3,
	53, 3, 53, 3, 53, 6, 53, 503, 10, 53, 13, 53, 14, 53, 504, 3, 53, 3, 53,
	3, 54, 3, 54, 3, 54, 5, 54, 512, 10, 54, 3, 55, 3, 55, 3, 55, 3, 55, 3,
	55, 3, 55, 3, 55, 3, 55, 3, 55, 5, 55, 523, 10, 55, 3, 55, 5, 55, 526,
	10, 55, 3, 56, 3, 56, 5, 56, 530, 10, 56, 3, 56, 5, 56, 533, 10, 56, 3,
	56, 3, 56, 3, 56, 6, 56, 538, 10, 56, 13, 56, 14, 56, 539, 3, 56, 3, 56,
	3, 57, 3, 57, 3, 58, 3, 58, 3, 58, 5, 58, 549, 10, 58, 3, 58, 5, 58, 552,
	10, 58, 3, 58, 3, 58, 3, 58, 3, 58, 6, 58, 558, 10, 58, 13, 58, 14, 58,
	559, 3, 58, 3, 58, 5, 58, 564, 10, 58, 5, 58, 566, 10, 58, 3, 59, 3, 59,
	5, 59, 570, 10, 59, 3, 59, 3, 59, 3, 59, 3, 59, 6, 59, 576, 10, 59, 13,
	59, 14, 59, 577, 3, 59, 3, 59, 3, 60, 3, 60, 3, 60, 5, 60, 585, 10, 60,
	3, 61, 3, 61, 3, 61, 3, 62, 3, 62, 3, 62, 3, 62, 3, 62, 6, 62, 595, 10,
	62, 13, 62, 14, 62, 596, 3, 62, 3, 62, 5, 62, 601, 10, 62, 3, 63, 3, 63,
	3, 63, 5, 63, 606, 10, 63, 3, 63, 3, 63, 3, 63, 3, 63, 6, 63, 612, 10,
	63, 13, 63, 14, 63, 613, 3, 63, 3, 63, 5, 63, 618, 10, 63, 3, 64, 3, 64,
	3, 64, 5, 64, 623, 10, 64, 3, 64, 3, 64, 3, 64, 3, 64, 6, 64, 629, 10,
	64, 13, 64, 14, 64, 630, 3, 64, 3, 64, 5, 64, 635, 10, 64, 3, 65, 3, 65,
	3, 65, 3, 65, 3, 65, 3, 65, 3, 65, 3, 65, 3, 65, 3, 65, 6, 65, 647, 10,
	65, 13, 65, 14, 65, 648, 3, 65, 3, 65, 3, 66, 7, 66, 654, 10, 66, 12, 66,
	14, 66, 657, 11, 66, 3, 66, 3, 66, 3, 66, 3, 66, 3, 67, 5, 67, 664, 10,
	67, 3, 67, 3, 67, 3, 67, 7, 67, 669, 10, 67, 12, 67, 14, 67, 672, 11, 67,
	3, 68, 3, 68, 7, 68, 676, 10, 68, 12, 68, 14, 68, 679, 11, 68, 3, 69, 6,
	69, 682, 10, 69, 13, 69, 14, 69, 683, 3, 70, 5, 70, 687, 10, 70, 3, 70,
	6, 70, 690, 10, 70, 13, 70, 14, 70, 691, 3, 70, 3, 70, 3, 70, 2, 2, 71,
	2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38,
	40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70, 72, 74,
	76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 100, 102, 104, 106, 108,
	110, 112, 114, 116, 118, 120, 122, 124, 126, 128, 130, 132, 134, 136, 138,
	2, 6, 4, 2, 5, 5, 57, 57, 3, 2, 8, 9, 3, 2, 56, 57, 4, 2, 57, 57, 61, 61,
	2, 726, 2, 140, 3, 2, 2, 2, 4, 149, 3, 2, 2, 2, 6, 157, 3, 2, 2, 2, 8,
	165, 3, 2, 2, 2, 10, 169, 3, 2, 2, 2, 12, 173, 3, 2, 2, 2, 14, 176, 3,
	2, 2, 2, 16, 178, 3, 2, 2, 2, 18, 189, 3, 2, 2, 2, 20, 193, 3, 2, 2, 2,
	22, 200, 3, 2, 2, 2, 24, 213, 3, 2, 2, 2, 26, 215, 3, 2, 2, 2, 28, 226,
	3, 2, 2, 2, 30, 231, 3, 2, 2, 2, 32, 233, 3, 2, 2, 2, 34, 235, 3, 2, 2,
	2, 36, 247, 3, 2, 2, 2, 38, 249, 3, 2, 2, 2, 40, 254, 3, 2, 2, 2, 42, 281,
	3, 2, 2, 2, 44, 283, 3, 2, 2, 2, 46, 290, 3, 2, 2, 2, 48, 315, 3, 2, 2,
	2, 50, 317, 3, 2, 2, 2, 52, 320, 3, 2, 2, 2, 54, 327, 3, 2, 2, 2, 56, 334,
	3, 2, 2, 2, 58, 337, 3, 2, 2, 2, 60, 349, 3, 2, 2, 2, 62, 357, 3, 2, 2,
	2, 64, 370, 3, 2, 2, 2, 66, 376, 3, 2, 2, 2, 68, 380, 3, 2, 2, 2, 70, 387,
	3, 2, 2, 2, 72, 396, 3, 2, 2, 2, 74, 402, 3, 2, 2, 2, 76, 404, 3, 2, 2,
	2, 78, 415, 3, 2, 2, 2, 80, 426, 3, 2, 2, 2, 82, 428, 3, 2, 2, 2, 84, 433,
	3, 2, 2, 2, 86, 435, 3, 2, 2, 2, 88, 437, 3, 2, 2, 2, 90, 441, 3, 2, 2,
	2, 92, 453, 3, 2, 2, 2, 94, 466, 3, 2, 2, 2, 96, 469, 3, 2, 2, 2, 98, 480,
	3, 2, 2, 2, 100, 485, 3, 2, 2, 2, 102, 488, 3, 2, 2, 2, 104, 498, 3, 2,
	2, 2, 106, 511, 3, 2, 2, 2, 108, 522, 3, 2, 2, 2, 110, 527, 3, 2, 2, 2,
	112, 543, 3, 2, 2, 2, 114, 565, 3, 2, 2, 2, 116, 567, 3, 2, 2, 2, 118,
	584, 3, 2, 2, 2, 120, 586, 3, 2, 2, 2, 122, 589, 3, 2, 2, 2, 124, 602,
	3, 2, 2, 2, 126, 619, 3, 2, 2, 2, 128, 636, 3, 2, 2, 2, 130, 655, 3, 2,
	2, 2, 132, 663, 3, 2, 2, 2, 134, 673, 3, 2, 2, 2, 136, 681, 3, 2, 2, 2,
	138, 686, 3, 2, 2, 2, 140, 141, 7, 26, 2, 2, 141, 146, 7, 57, 2, 2, 142,
	143, 7, 25, 2, 2, 143, 145, 7, 57, 2, 2, 144, 142, 3, 2, 2, 2, 145, 148,
	3, 2, 2, 2, 146, 144, 3, 2, 2, 2, 146, 147, 3, 2, 2, 2, 147, 3, 3, 2, 2,
	2, 148, 146, 3, 2, 2, 2, 149, 150, 7, 44, 2, 2, 150, 153, 7, 52, 2, 2,
	151, 152, 7, 35, 2, 2, 152, 154, 7, 52, 2, 2, 153, 151, 3, 2, 2, 2, 153,
	154, 3, 2, 2, 2, 154, 155, 3, 2, 2, 2, 155, 156, 7, 45, 2, 2, 156, 5, 3,
	2, 2, 2, 157, 162, 5, 2, 2, 2, 158, 159, 7, 27, 2, 2, 159, 161, 5, 2, 2,
	2, 160, 158, 3, 2, 2, 2, 161, 164, 3, 2, 2, 2, 162, 160, 3, 2, 2, 2, 162,
	163, 3, 2, 2, 2, 163, 7, 3, 2, 2, 2, 164, 162, 3, 2, 2, 2, 165, 166, 7,
	40, 2, 2, 166, 167, 5, 6, 4, 2, 167, 168, 7, 41, 2, 2, 168, 9, 3, 2, 2,
	2, 169, 170, 7, 57, 2, 2, 170, 171, 7, 35, 2, 2, 171, 172, 7, 57, 2, 2,
	172, 11, 3, 2, 2, 2, 173, 174, 7, 47, 2, 2, 174, 175, 7, 59, 2, 2, 175,
	13, 3, 2, 2, 2, 176, 177, 7, 53, 2, 2, 177, 15, 3, 2, 2, 2, 178, 179, 7,
	40, 2, 2, 179, 184, 5, 14, 8, 2, 180, 181, 7, 27, 2, 2, 181, 183, 5, 14,
	8, 2, 182, 180, 3, 2, 2, 2, 183, 186, 3, 2, 2, 2, 184, 182, 3, 2, 2, 2,
	184, 185, 3, 2, 2, 2, 185, 187, 3, 2, 2, 2, 186, 184, 3, 2, 2, 2, 187,
	188, 7, 41, 2, 2, 188, 17, 3, 2, 2, 2, 189, 190, 7, 40, 2, 2, 190, 191,
	5, 16, 9, 2, 191, 192, 7, 41, 2, 2, 192, 19, 3, 2, 2, 2, 193, 194, 7, 57,
	2, 2, 194, 198, 7, 28, 2, 2, 195, 199, 5, 14, 8, 2, 196, 199, 5, 16, 9,
	2, 197, 199, 5, 18, 10, 2, 198, 195, 3, 2, 2, 2, 198, 196, 3, 2, 2, 2,
	198, 197, 3, 2, 2, 2, 199, 21, 3, 2, 2, 2, 200, 201, 7, 40, 2, 2, 201,
	206, 5, 20, 11, 2, 202, 203, 7, 27, 2, 2, 203, 205, 5, 20, 11, 2, 204,
	202, 3, 2, 2, 2, 205, 208, 3, 2, 2, 2, 206, 204, 3, 2, 2, 2, 206, 207,
	3, 2, 2, 2, 207, 209, 3, 2, 2, 2, 208, 206, 3, 2, 2, 2, 209, 210, 7, 41,
	2, 2, 210, 23, 3, 2, 2, 2, 211, 214, 5, 20, 11, 2, 212, 214, 5, 2, 2, 2,
	213, 211, 3, 2, 2, 2, 213, 212, 3, 2, 2, 2, 214, 25, 3, 2, 2, 2, 215, 216,
	7, 40, 2, 2, 216, 221, 5, 24, 13, 2, 217, 218, 7, 27, 2, 2, 218, 220, 5,
	24, 13, 2, 219, 217, 3, 2, 2, 2, 220, 223, 3, 2, 2, 2, 221, 219, 3, 2,
	2, 2, 221, 222, 3, 2, 2, 2, 222, 224, 3, 2, 2, 2, 223, 221, 3, 2, 2, 2,
	224, 225, 7, 41, 2, 2, 225, 27, 3, 2, 2, 2, 226, 227, 7, 17, 2, 2, 227,
	229, 9, 2, 2, 2, 228, 230, 5, 4, 3, 2, 229, 228, 3, 2, 2, 2, 229, 230,
	3, 2, 2, 2, 230, 29, 3, 2, 2, 2, 231, 232, 5, 28, 15, 2, 232, 31, 3, 2,
	2, 2, 233, 234, 7, 57, 2, 2, 234, 33, 3, 2, 2, 2, 235, 236, 7, 33, 2, 2,
	236, 238, 7, 3, 2, 2, 237, 239, 5, 12, 7, 2, 238, 237, 3, 2, 2, 2, 239,
	240, 3, 2, 2, 2, 240, 238, 3, 2, 2, 2, 240, 241, 3, 2, 2, 2, 241, 242,
	3, 2, 2, 2, 242, 243, 7, 4, 2, 2, 243, 35, 3, 2, 2, 2, 244, 248, 7, 53,
	2, 2, 245, 248, 5, 16, 9, 2, 246, 248, 5, 34, 18, 2, 247, 244, 3, 2, 2,
	2, 247, 245, 3, 2, 2, 2, 247, 246, 3, 2, 2, 2, 248, 37, 3, 2, 2, 2, 249,
	250, 7, 38, 2, 2, 250, 251, 7, 63, 2, 2, 251, 252, 7, 28, 2, 2, 252, 253,
	5, 36, 19, 2, 253, 39, 3, 2, 2, 2, 254, 256, 7, 3, 2, 2, 255, 257, 5, 38,
	20, 2, 256, 255, 3, 2, 2, 2, 257, 258, 3, 2, 2, 2, 258, 256, 3, 2, 2, 2,
	258, 259, 3, 2, 2, 2, 259, 260, 3, 2, 2, 2, 260, 261, 7, 4, 2, 2, 261,
	41, 3, 2, 2, 2, 262, 282, 5, 30, 16, 2, 263, 267, 5, 10, 6, 2, 264, 267,
	7, 5, 2, 2, 265, 267, 5, 32, 17, 2, 266, 263, 3, 2, 2, 2, 266, 264, 3,
	2, 2, 2, 266, 265, 3, 2, 2, 2, 267, 269, 3, 2, 2, 2, 268, 270, 5, 4, 3,
	2, 269, 268, 3, 2, 2, 2, 269, 270, 3, 2, 2, 2, 270, 272, 3, 2, 2, 2, 271,
	273, 7, 37, 2, 2, 272, 271, 3, 2, 2, 2, 272, 273, 3, 2, 2, 2, 273, 275,
	3, 2, 2, 2, 274, 276, 5, 26, 14, 2, 275, 274, 3, 2, 2, 2, 275, 276, 3,
	2, 2, 2, 276, 279, 3, 2, 2, 2, 277, 278, 7, 33, 2, 2, 278, 280, 5, 40,
	21, 2, 279, 277, 3, 2, 2, 2, 279, 280, 3, 2, 2, 2, 280, 282, 3, 2, 2, 2,
	281, 262, 3, 2, 2, 2, 281, 266, 3, 2, 2, 2, 282, 43, 3, 2, 2, 2, 283, 284,
	7, 57, 2, 2, 284, 285, 7, 21, 2, 2, 285, 286, 5, 42, 22, 2, 286, 45, 3,
	2, 2, 2, 287, 289, 7, 55, 2, 2, 288, 287, 3, 2, 2, 2, 289, 292, 3, 2, 2,
	2, 290, 288, 3, 2, 2, 2, 290, 291, 3, 2, 2, 2, 291, 293, 3, 2, 2, 2, 292,
	290, 3, 2, 2, 2, 293, 294, 9, 3, 2, 2, 294, 296, 7, 57, 2, 2, 295, 297,
	5, 26, 14, 2, 296, 295, 3, 2, 2, 2, 296, 297, 3, 2, 2, 2, 297, 298, 3,
	2, 2, 2, 298, 299, 7, 33, 2, 2, 299, 302, 7, 3, 2, 2, 300, 303, 7, 55,
	2, 2, 301, 303, 5, 44, 23, 2, 302, 300, 3, 2, 2, 2, 302, 301, 3, 2, 2,
	2, 303, 304, 3, 2, 2, 2, 304, 302, 3, 2, 2, 2, 304, 305, 3, 2, 2, 2, 305,
	306, 3, 2, 2, 2, 306, 307, 7, 4, 2, 2, 307, 47, 3, 2, 2, 2, 308, 311, 7,
	57, 2, 2, 309, 310, 7, 35, 2, 2, 310, 312, 7, 57, 2, 2, 311, 309, 3, 2,
	2, 2, 311, 312, 3, 2, 2, 2, 312, 316, 3, 2, 2, 2, 313, 316, 7, 56, 2, 2,
	314, 316, 7, 61, 2, 2, 315, 308, 3, 2, 2, 2, 315, 313, 3, 2, 2, 2, 315,
	314, 3, 2, 2, 2, 316, 49, 3, 2, 2, 2, 317, 318, 7, 20, 2, 2, 318, 319,
	5, 48, 25, 2, 319, 51, 3, 2, 2, 2, 320, 324, 5, 48, 25, 2, 321, 323, 5,
	50, 26, 2, 322, 321, 3, 2, 2, 2, 323, 326, 3, 2, 2, 2, 324, 322, 3, 2,
	2, 2, 324, 325, 3, 2, 2, 2, 325, 53, 3, 2, 2, 2, 326, 324, 3, 2, 2, 2,
	327, 329, 5, 52, 27, 2, 328, 330, 7, 53, 2, 2, 329, 328, 3, 2, 2, 2, 329,
	330, 3, 2, 2, 2, 330, 332, 3, 2, 2, 2, 331, 333, 5, 26, 14, 2, 332, 331,
	3, 2, 2, 2, 332, 333, 3, 2, 2, 2, 333, 55, 3, 2, 2, 2, 334, 335, 7, 57,
	2, 2, 335, 336, 7, 33, 2, 2, 336, 57, 3, 2, 2, 2, 337, 338, 7, 33, 2, 2,
	338, 343, 7, 3, 2, 2, 339, 341, 7, 57, 2, 2, 340, 342, 5, 26, 14, 2, 341,
	340, 3, 2, 2, 2, 341, 342, 3, 2, 2, 2, 342, 344, 3, 2, 2, 2, 343, 339,
	3, 2, 2, 2, 344, 345, 3, 2, 2, 2, 345, 343, 3, 2, 2, 2, 345, 346, 3, 2,
	2, 2, 346, 347, 3, 2, 2, 2, 347, 348, 7, 4, 2, 2, 348, 59, 3, 2, 2, 2,
	349, 350, 9, 3, 2, 2, 350, 352, 7, 57, 2, 2, 351, 353, 5, 58, 30, 2, 352,
	351, 3, 2, 2, 2, 352, 353, 3, 2, 2, 2, 353, 61, 3, 2, 2, 2, 354, 356, 7,
	55, 2, 2, 355, 354, 3, 2, 2, 2, 356, 359, 3, 2, 2, 2, 357, 355, 3, 2, 2,
	2, 357, 358, 3, 2, 2, 2, 358, 360, 3, 2, 2, 2, 359, 357, 3, 2, 2, 2, 360,
	361, 7, 7, 2, 2, 361, 362, 5, 56, 29, 2, 362, 364, 7, 3, 2, 2, 363, 365,
	5, 60, 31, 2, 364, 363, 3, 2, 2, 2, 365, 366, 3, 2, 2, 2, 366, 364, 3,
	2, 2, 2, 366, 367, 3, 2, 2, 2, 367, 368, 3, 2, 2, 2, 368, 369, 7, 4, 2,
	2, 369, 63, 3, 2, 2, 2, 370, 371, 7, 38, 2, 2, 371, 372, 7, 57, 2, 2, 372,
	373, 7, 28, 2, 2, 373, 374, 7, 53, 2, 2, 374, 375, 7, 54, 2, 2, 375, 65,
	3, 2, 2, 2, 376, 377, 7, 42, 2, 2, 377, 378, 7, 57, 2, 2, 378, 379, 7,
	43, 2, 2, 379, 67, 3, 2, 2, 2, 380, 381, 7, 57, 2, 2, 381, 385, 7, 28,
	2, 2, 382, 386, 7, 57, 2, 2, 383, 386, 7, 5, 2, 2, 384, 386, 5, 66, 34,
	2, 385, 382, 3, 2, 2, 2, 385, 383, 3, 2, 2, 2, 385, 384, 3, 2, 2, 2, 386,
	69, 3, 2, 2, 2, 387, 388, 7, 37, 2, 2, 388, 393, 5, 68, 35, 2, 389, 390,
	7, 39, 2, 2, 390, 392, 5, 68, 35, 2, 391, 389, 3, 2, 2, 2, 392, 395, 3,
	2, 2, 2, 393, 391, 3, 2, 2, 2, 393, 394, 3, 2, 2, 2, 394, 71, 3, 2, 2,
	2, 395, 393, 3, 2, 2, 2, 396, 397, 7, 42, 2, 2, 397, 398, 7, 57, 2, 2,
	398, 399, 7, 21, 2, 2, 399, 400, 9, 2, 2, 2, 400, 401, 7, 43, 2, 2, 401,
	73, 3, 2, 2, 2, 402, 403, 9, 4, 2, 2, 403, 75, 3, 2, 2, 2, 404, 407, 7,
	30, 2, 2, 405, 408, 5, 74, 38, 2, 406, 408, 5, 72, 37, 2, 407, 405, 3,
	2, 2, 2, 407, 406, 3, 2, 2, 2, 408, 77, 3, 2, 2, 2, 409, 416, 7, 30, 2,
	2, 410, 412, 5, 76, 39, 2, 411, 410, 3, 2, 2, 2, 412, 413, 3, 2, 2, 2,
	413, 411, 3, 2, 2, 2, 413, 414, 3, 2, 2, 2, 414, 416, 3, 2, 2, 2, 415,
	409, 3, 2, 2, 2, 415, 411, 3, 2, 2, 2, 416, 79, 3, 2, 2, 2, 417, 419, 7,
	35, 2, 2, 418, 417, 3, 2, 2, 2, 418, 419, 3, 2, 2, 2, 419, 421, 3, 2, 2,
	2, 420, 422, 7, 57, 2, 2, 421, 420, 3, 2, 2, 2, 422, 423, 3, 2, 2, 2, 423,
	421, 3, 2, 2, 2, 423, 424, 3, 2, 2, 2, 424, 427, 3, 2, 2, 2, 425, 427,
	7, 56, 2, 2, 426, 418, 3, 2, 2, 2, 426, 425, 3, 2, 2, 2, 427, 81, 3, 2,
	2, 2, 428, 429, 7, 11, 2, 2, 429, 430, 7, 59, 2, 2, 430, 83, 3, 2, 2, 2,
	431, 434, 7, 35, 2, 2, 432, 434, 5, 52, 27, 2, 433, 431, 3, 2, 2, 2, 433,
	432, 3, 2, 2, 2, 434, 85, 3, 2, 2, 2, 435, 436, 9, 5, 2, 2, 436, 87, 3,
	2, 2, 2, 437, 438, 5, 84, 43, 2, 438, 439, 7, 22, 2, 2, 439, 440, 5, 86,
	44, 2, 440, 89, 3, 2, 2, 2, 441, 442, 7, 12, 2, 2, 442, 443, 7, 61, 2,
	2, 443, 444, 7, 33, 2, 2, 444, 448, 7, 3, 2, 2, 445, 447, 5, 108, 55, 2,
	446, 445, 3, 2, 2, 2, 447, 450, 3, 2, 2, 2, 448, 446, 3, 2, 2, 2, 448,
	449, 3, 2, 2, 2, 449, 451, 3, 2, 2, 2, 450, 448, 3, 2, 2, 2, 451, 452,
	7, 4, 2, 2, 452, 91, 3, 2, 2, 2, 453, 464, 5, 90, 46, 2, 454, 455, 7, 13,
	2, 2, 455, 456, 7, 33, 2, 2, 456, 460, 7, 3, 2, 2, 457, 459, 5, 108, 55,
	2, 458, 457, 3, 2, 2, 2, 459, 462, 3, 2, 2, 2, 460, 458, 3, 2, 2, 2, 460,
	461, 3, 2, 2, 2, 461, 463, 3, 2, 2, 2, 462, 460, 3, 2, 2, 2, 463, 465,
	7, 4, 2, 2, 464, 454, 3, 2, 2, 2, 464, 465, 3, 2, 2, 2, 465, 93, 3, 2,
	2, 2, 466, 467, 7, 61, 2, 2, 467, 468, 7, 33, 2, 2, 468, 95, 3, 2, 2, 2,
	469, 470, 7, 14, 2, 2, 470, 471, 5, 94, 48, 2, 471, 475, 7, 3, 2, 2, 472,
	474, 5, 108, 55, 2, 473, 472, 3, 2, 2, 2, 474, 477, 3, 2, 2, 2, 475, 473,
	3, 2, 2, 2, 475, 476, 3, 2, 2, 2, 476, 478, 3, 2, 2, 2, 477, 475, 3, 2,
	2, 2, 478, 479, 7, 4, 2, 2, 479, 97, 3, 2, 2, 2, 480, 481, 7, 55, 2, 2,
	481, 99, 3, 2, 2, 2, 482, 484, 7, 57, 2, 2, 483, 482, 3, 2, 2, 2, 484,
	487, 3, 2, 2, 2, 485, 483, 3, 2, 2, 2, 485, 486, 3, 2, 2, 2, 486, 101,
	3, 2, 2, 2, 487, 485, 3, 2, 2, 2, 488, 489, 5, 100, 51, 2, 489, 490, 7,
	33, 2, 2, 490, 492, 7, 3, 2, 2, 491, 493, 5, 108, 55, 2, 492, 491, 3, 2,
	2, 2, 493, 494, 3, 2, 2, 2, 494, 492, 3, 2, 2, 2, 494, 495, 3, 2, 2, 2,
	495, 496, 3, 2, 2, 2, 496, 497, 7, 4, 2, 2, 497, 103, 3, 2, 2, 2, 498,
	499, 7, 18, 2, 2, 499, 500, 7, 33, 2, 2, 500, 502, 7, 3, 2, 2, 501, 503,
	5, 102, 52, 2, 502, 501, 3, 2, 2, 2, 503, 504, 3, 2, 2, 2, 504, 502, 3,
	2, 2, 2, 504, 505, 3, 2, 2, 2, 505, 506, 3, 2, 2, 2, 506, 507, 7, 4, 2,
	2, 507, 105, 3, 2, 2, 2, 508, 512, 5, 12, 7, 2, 509, 512, 7, 56, 2, 2,
	510, 512, 7, 16, 2, 2, 511, 508, 3, 2, 2, 2, 511, 509, 3, 2, 2, 2, 511,
	510, 3, 2, 2, 2, 512, 107, 3, 2, 2, 2, 513, 523, 5, 92, 47, 2, 514, 523,
	5, 96, 49, 2, 515, 523, 5, 82, 42, 2, 516, 523, 5, 88, 45, 2, 517, 523,
	5, 104, 53, 2, 518, 523, 5, 98, 50, 2, 519, 523, 7, 53, 2, 2, 520, 523,
	5, 106, 54, 2, 521, 523, 5, 38, 20, 2, 522, 513, 3, 2, 2, 2, 522, 514,
	3, 2, 2, 2, 522, 515, 3, 2, 2, 2, 522, 516, 3, 2, 2, 2, 522, 517, 3, 2,
	2, 2, 522, 518, 3, 2, 2, 2, 522, 519, 3, 2, 2, 2, 522, 520, 3, 2, 2, 2,
	522, 521, 3, 2, 2, 2, 523, 525, 3, 2, 2, 2, 524, 526, 5, 26, 14, 2, 525,
	524, 3, 2, 2, 2, 525, 526, 3, 2, 2, 2, 526, 109, 3, 2, 2, 2, 527, 529,
	7, 6, 2, 2, 528, 530, 5, 70, 36, 2, 529, 528, 3, 2, 2, 2, 529, 530, 3,
	2, 2, 2, 530, 532, 3, 2, 2, 2, 531, 533, 5, 22, 12, 2, 532, 531, 3, 2,
	2, 2, 532, 533, 3, 2, 2, 2, 533, 534, 3, 2, 2, 2, 534, 535, 7, 33, 2, 2,
	535, 537, 7, 3, 2, 2, 536, 538, 5, 108, 55, 2, 537, 536, 3, 2, 2, 2, 538,
	539, 3, 2, 2, 2, 539, 537, 3, 2, 2, 2, 539, 540, 3, 2, 2, 2, 540, 541,
	3, 2, 2, 2, 541, 542, 7, 4, 2, 2, 542, 111, 3, 2, 2, 2, 543, 544, 7, 16,
	2, 2, 544, 113, 3, 2, 2, 2, 545, 566, 7, 16, 2, 2, 546, 548, 5, 80, 41,
	2, 547, 549, 7, 53, 2, 2, 548, 547, 3, 2, 2, 2, 548, 549, 3, 2, 2, 2, 549,
	551, 3, 2, 2, 2, 550, 552, 5, 26, 14, 2, 551, 550, 3, 2, 2, 2, 551, 552,
	3, 2, 2, 2, 552, 553, 3, 2, 2, 2, 553, 563, 7, 33, 2, 2, 554, 564, 5, 112,
	57, 2, 555, 557, 7, 3, 2, 2, 556, 558, 5, 108, 55, 2, 557, 556, 3, 2, 2,
	2, 558, 559, 3, 2, 2, 2, 559, 557, 3, 2, 2, 2, 559, 560, 3, 2, 2, 2, 560,
	561, 3, 2, 2, 2, 561, 562, 7, 4, 2, 2, 562, 564, 3, 2, 2, 2, 563, 554,
	3, 2, 2, 2, 563, 555, 3, 2, 2, 2, 564, 566, 3, 2, 2, 2, 565, 545, 3, 2,
	2, 2, 565, 546, 3, 2, 2, 2, 566, 115, 3, 2, 2, 2, 567, 569, 5, 78, 40,
	2, 568, 570, 5, 26, 14, 2, 569, 568, 3, 2, 2, 2, 569, 570, 3, 2, 2, 2,
	570, 571, 3, 2, 2, 2, 571, 572, 7, 33, 2, 2, 572, 575, 7, 3, 2, 2, 573,
	576, 5, 110, 56, 2, 574, 576, 5, 116, 59, 2, 575, 573, 3, 2, 2, 2, 575,
	574, 3, 2, 2, 2, 576, 577, 3, 2, 2, 2, 577, 575, 3, 2, 2, 2, 577, 578,
	3, 2, 2, 2, 578, 579, 3, 2, 2, 2, 579, 580, 7, 4, 2, 2, 580, 117, 3, 2,
	2, 2, 581, 585, 5, 88, 45, 2, 582, 583, 7, 6, 2, 2, 583, 585, 5, 78, 40,
	2, 584, 581, 3, 2, 2, 2, 584, 582, 3, 2, 2, 2, 585, 119, 3, 2, 2, 2, 586,
	587, 5, 118, 60, 2, 587, 588, 5, 26, 14, 2, 588, 121, 3, 2, 2, 2, 589,
	590, 7, 24, 2, 2, 590, 600, 7, 33, 2, 2, 591, 601, 7, 16, 2, 2, 592, 594,
	7, 3, 2, 2, 593, 595, 5, 120, 61, 2, 594, 593, 3, 2, 2, 2, 595, 596, 3,
	2, 2, 2, 596, 594, 3, 2, 2, 2, 596, 597, 3, 2, 2, 2, 597, 598, 3, 2, 2,
	2, 598, 599, 7, 4, 2, 2, 599, 601, 3, 2, 2, 2, 600, 591, 3, 2, 2, 2, 600,
	592, 3, 2, 2, 2, 601, 123, 3, 2, 2, 2, 602, 603, 7, 19, 2, 2, 603, 605,
	7, 65, 2, 2, 604, 606, 5, 26, 14, 2, 605, 604, 3, 2, 2, 2, 605, 606, 3,
	2, 2, 2, 606, 607, 3, 2, 2, 2, 607, 617, 7, 33, 2, 2, 608, 618, 7, 16,
	2, 2, 609, 611, 7, 3, 2, 2, 610, 612, 5, 108, 55, 2, 611, 610, 3, 2, 2,
	2, 612, 613, 3, 2, 2, 2, 613, 611, 3, 2, 2, 2, 613, 614, 3, 2, 2, 2, 614,
	615, 3, 2, 2, 2, 615, 616, 7, 4, 2, 2, 616, 618, 3, 2, 2, 2, 617, 608,
	3, 2, 2, 2, 617, 609, 3, 2, 2, 2, 618, 125, 3, 2, 2, 2, 619, 620, 7, 57,
	2, 2, 620, 622, 7, 23, 2, 2, 621, 623, 5, 26, 14, 2, 622, 621, 3, 2, 2,
	2, 622, 623, 3, 2, 2, 2, 623, 624, 3, 2, 2, 2, 624, 634, 7, 33, 2, 2, 625,
	635, 7, 16, 2, 2, 626, 628, 7, 3, 2, 2, 627, 629, 5, 108, 55, 2, 628, 627,
	3, 2, 2, 2, 629, 630, 3, 2, 2, 2, 630, 628, 3, 2, 2, 2, 630, 631, 3, 2,
	2, 2, 631, 632, 3, 2, 2, 2, 632, 633, 7, 4, 2, 2, 633, 635, 3, 2, 2, 2,
	634, 625, 3, 2, 2, 2, 634, 626, 3, 2, 2, 2, 635, 127, 3, 2, 2, 2, 636,
	646, 7, 3, 2, 2, 637, 647, 5, 46, 24, 2, 638, 647, 5, 62, 32, 2, 639, 647,
	7, 55, 2, 2, 640, 647, 5, 116, 59, 2, 641, 647, 5, 114, 58, 2, 642, 647,
	5, 122, 62, 2, 643, 647, 5, 124, 63, 2, 644, 647, 5, 126, 64, 2, 645, 647,
	5, 38, 20, 2, 646, 637, 3, 2, 2, 2, 646, 638, 3, 2, 2, 2, 646, 639, 3,
	2, 2, 2, 646, 640, 3, 2, 2, 2, 646, 641, 3, 2, 2, 2, 646, 642, 3, 2, 2,
	2, 646, 643, 3, 2, 2, 2, 646, 644, 3, 2, 2, 2, 646, 645, 3, 2, 2, 2, 647,
	648, 3, 2, 2, 2, 648, 646, 3, 2, 2, 2, 648, 649, 3, 2, 2, 2, 649, 650,
	3, 2, 2, 2, 650, 651, 7, 4, 2, 2, 651, 129, 3, 2, 2, 2, 652, 654, 7, 55,
	2, 2, 653, 652, 3, 2, 2, 2, 654, 657, 3, 2, 2, 2, 655, 653, 3, 2, 2, 2,
	655, 656, 3, 2, 2, 2, 656, 658, 3, 2, 2, 2, 657, 655, 3, 2, 2, 2, 658,
	659, 5, 54, 28, 2, 659, 660, 7, 33, 2, 2, 660, 661, 5, 128, 65, 2, 661,
	131, 3, 2, 2, 2, 662, 664, 7, 30, 2, 2, 663, 662, 3, 2, 2, 2, 663, 664,
	3, 2, 2, 2, 664, 665, 3, 2, 2, 2, 665, 670, 7, 57, 2, 2, 666, 667, 7, 30,
	2, 2, 667, 669, 7, 57, 2, 2, 668, 666, 3, 2, 2, 2, 669, 672, 3, 2, 2, 2,
	670, 668, 3, 2, 2, 2, 670, 671, 3, 2, 2, 2, 671, 133, 3, 2, 2, 2, 672,
	670, 3, 2, 2, 2, 673, 677, 7, 10, 2, 2, 674, 676, 7, 55, 2, 2, 675, 674,
	3, 2, 2, 2, 676, 679, 3, 2, 2, 2, 677, 675, 3, 2, 2, 2, 677, 678, 3, 2,
	2, 2, 678, 135, 3, 2, 2, 2, 679, 677, 3, 2, 2, 2, 680, 682, 5, 134, 68,
	2, 681, 680, 3, 2, 2, 2, 682, 683, 3, 2, 2, 2, 683, 681, 3, 2, 2, 2, 683,
	684, 3, 2, 2, 2, 684, 137, 3, 2, 2, 2, 685, 687, 5, 136, 69, 2, 686, 685,
	3, 2, 2, 2, 686, 687, 3, 2, 2, 2, 687, 689, 3, 2, 2, 2, 688, 690, 5, 130,
	66, 2, 689, 688, 3, 2, 2, 2, 690, 691, 3, 2, 2, 2, 691, 689, 3, 2, 2, 2,
	691, 692, 3, 2, 2, 2, 692, 693, 3, 2, 2, 2, 693, 694, 7, 2, 2, 3, 694,
	139, 3, 2, 2, 2, 82, 146, 153, 162, 184, 198, 206, 213, 221, 229, 240,
	247, 258, 266, 269, 272, 275, 279, 281, 290, 296, 302, 304, 311, 315, 324,
	329, 332, 341, 345, 352, 357, 366, 385, 393, 407, 413, 415, 418, 423, 426,
	433, 448, 460, 464, 475, 485, 494, 504, 511, 522, 525, 529, 532, 539, 548,
	551, 559, 563, 565, 569, 575, 577, 584, 596, 600, 605, 613, 617, 622, 630,
	634, 646, 648, 655, 663, 670, 677, 683, 686, 691,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "", "", "", "", "'!wrap'", "'!table'", "'!type'", "", "", "", "", "",
	"", "'...'", "'set of'", "", "'<->'", "'::'", "'<:'", "'<-'", "'->'", "'.. * <- *'",
	"'+'", "'~'", "','", "'='", "'$'", "'/'", "'-'", "'*'", "':'", "'%'", "'.'",
	"'!'", "'?'", "'@'", "'&'", "'['", "']'", "'{'", "'}'", "'('", "')'", "'#'",
	"'|'",
}
var symbolicNames = []string{
	"", "INDENT", "DEDENT", "NativeDataTypes", "HTTP_VERBS", "WRAP", "TABLE",
	"TYPE", "IMPORT", "RETURN", "IF", "ELSE", "FOR", "LOOP", "WHATEVER", "SET_OF",
	"ONE_OF", "DISTANCE", "NAME_SEP", "LESS_COLON", "ARROW_LEFT", "ARROW_RIGHT",
	"COLLECTOR", "PLUS", "TILDE", "COMMA", "EQ", "DOLLAR", "FORWARD_SLASH",
	"MINUS", "STAR", "COLON", "PERCENT", "DOT", "EXCLAIM", "QN", "AT", "AMP",
	"SQ_OPEN", "SQ_CLOSE", "CURLY_OPEN", "CURLY_CLOSE", "OPEN_PAREN", "CLOSE_PAREN",
	"HASH", "PIPE", "DBL_QT", "SINGLE_QT", "EMPTY_LINE", "INDENTED_COMMENT",
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
	"documentation_stmts", "var_in_curly", "query_var", "query_param", "http_path_var_with_type",
	"http_path_static", "http_path_suffix", "http_path", "endpoint_name", "ret_stmt",
	"target", "target_endpoint", "call_stmt", "if_stmt", "if_else", "for_cond",
	"for_stmt", "http_method_comment", "one_of_case_label", "one_of_cases",
	"one_of_stmt", "text_stmt", "statements", "method_def", "shortcut", "simple_endpoint",
	"rest_endpoint", "collector_stmt", "collector_stmts", "collector", "event",
	"subscribe", "app_decl", "application", "path", "import_stmt", "imports_decl",
	"sysl_file",
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
	SyslParserWHATEVER         = 14
	SyslParserSET_OF           = 15
	SyslParserONE_OF           = 16
	SyslParserDISTANCE         = 17
	SyslParserNAME_SEP         = 18
	SyslParserLESS_COLON       = 19
	SyslParserARROW_LEFT       = 20
	SyslParserARROW_RIGHT      = 21
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
	SyslParserHASH             = 44
	SyslParserPIPE             = 45
	SyslParserDBL_QT           = 46
	SyslParserSINGLE_QT        = 47
	SyslParserEMPTY_LINE       = 48
	SyslParserINDENTED_COMMENT = 49
	SyslParserDIGITS           = 50
	SyslParserQSTRING          = 51
	SyslParserNEWLINE          = 52
	SyslParserSYSL_COMMENT     = 53
	SyslParserTEXT_LINE        = 54
	SyslParserName             = 55
	SyslParserWS               = 56
	SyslParserTEXT             = 57
	SyslParserSKIP_WS          = 58
	SyslParserTEXT_NAME        = 59
	SyslParserPOP_WS           = 60
	SyslParserVAR_NAME         = 61
	SyslParserSKIP_WS_2        = 62
	SyslParserEVENT_NAME       = 63
)

// SyslParser rules.
const (
	SyslParserRULE_modifier                = 0
	SyslParserRULE_size_spec               = 1
	SyslParserRULE_modifier_list           = 2
	SyslParserRULE_modifiers               = 3
	SyslParserRULE_reference               = 4
	SyslParserRULE_doc_string              = 5
	SyslParserRULE_quoted_string           = 6
	SyslParserRULE_array_of_strings        = 7
	SyslParserRULE_array_of_arrays         = 8
	SyslParserRULE_nvp                     = 9
	SyslParserRULE_attributes              = 10
	SyslParserRULE_entry                   = 11
	SyslParserRULE_attribs_or_modifiers    = 12
	SyslParserRULE_set_type                = 13
	SyslParserRULE_collection_type         = 14
	SyslParserRULE_user_defined_type       = 15
	SyslParserRULE_multi_line_docstring    = 16
	SyslParserRULE_annotation_value        = 17
	SyslParserRULE_annotation              = 18
	SyslParserRULE_annotations             = 19
	SyslParserRULE_field_type              = 20
	SyslParserRULE_field                   = 21
	SyslParserRULE_table                   = 22
	SyslParserRULE_package_name            = 23
	SyslParserRULE_sub_package             = 24
	SyslParserRULE_app_name                = 25
	SyslParserRULE_name_with_attribs       = 26
	SyslParserRULE_model_name              = 27
	SyslParserRULE_inplace_table_def       = 28
	SyslParserRULE_table_refs              = 29
	SyslParserRULE_facade                  = 30
	SyslParserRULE_documentation_stmts     = 31
	SyslParserRULE_var_in_curly            = 32
	SyslParserRULE_query_var               = 33
	SyslParserRULE_query_param             = 34
	SyslParserRULE_http_path_var_with_type = 35
	SyslParserRULE_http_path_static        = 36
	SyslParserRULE_http_path_suffix        = 37
	SyslParserRULE_http_path               = 38
	SyslParserRULE_endpoint_name           = 39
	SyslParserRULE_ret_stmt                = 40
	SyslParserRULE_target                  = 41
	SyslParserRULE_target_endpoint         = 42
	SyslParserRULE_call_stmt               = 43
	SyslParserRULE_if_stmt                 = 44
	SyslParserRULE_if_else                 = 45
	SyslParserRULE_for_cond                = 46
	SyslParserRULE_for_stmt                = 47
	SyslParserRULE_http_method_comment     = 48
	SyslParserRULE_one_of_case_label       = 49
	SyslParserRULE_one_of_cases            = 50
	SyslParserRULE_one_of_stmt             = 51
	SyslParserRULE_text_stmt               = 52
	SyslParserRULE_statements              = 53
	SyslParserRULE_method_def              = 54
	SyslParserRULE_shortcut                = 55
	SyslParserRULE_simple_endpoint         = 56
	SyslParserRULE_rest_endpoint           = 57
	SyslParserRULE_collector_stmt          = 58
	SyslParserRULE_collector_stmts         = 59
	SyslParserRULE_collector               = 60
	SyslParserRULE_event                   = 61
	SyslParserRULE_subscribe               = 62
	SyslParserRULE_app_decl                = 63
	SyslParserRULE_application             = 64
	SyslParserRULE_path                    = 65
	SyslParserRULE_import_stmt             = 66
	SyslParserRULE_imports_decl            = 67
	SyslParserRULE_sysl_file               = 68
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
		p.SetState(138)
		p.Match(SyslParserTILDE)
	}
	{
		p.SetState(139)
		p.Match(SyslParserName)
	}
	p.SetState(144)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SyslParserPLUS {
		{
			p.SetState(140)
			p.Match(SyslParserPLUS)
		}
		{
			p.SetState(141)
			p.Match(SyslParserName)
		}

		p.SetState(146)
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
		p.SetState(147)
		p.Match(SyslParserOPEN_PAREN)
	}
	{
		p.SetState(148)
		p.Match(SyslParserDIGITS)
	}
	p.SetState(151)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserDOT {
		{
			p.SetState(149)
			p.Match(SyslParserDOT)
		}
		{
			p.SetState(150)
			p.Match(SyslParserDIGITS)
		}

	}
	{
		p.SetState(153)
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
		p.SetState(155)
		p.Modifier()
	}
	p.SetState(160)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SyslParserCOMMA {
		{
			p.SetState(156)
			p.Match(SyslParserCOMMA)
		}
		{
			p.SetState(157)
			p.Modifier()
		}

		p.SetState(162)
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
		p.SetState(163)
		p.Match(SyslParserSQ_OPEN)
	}
	{
		p.SetState(164)
		p.Modifier_list()
	}
	{
		p.SetState(165)
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
		p.SetState(167)

		var _m = p.Match(SyslParserName)

		localctx.(*ReferenceContext).parent_ref = _m
	}
	{
		p.SetState(168)
		p.Match(SyslParserDOT)
	}
	{
		p.SetState(169)

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
		p.SetState(171)
		p.Match(SyslParserPIPE)
	}
	{
		p.SetState(172)
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
		p.SetState(174)
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
		p.SetState(176)
		p.Match(SyslParserSQ_OPEN)
	}
	{
		p.SetState(177)
		p.Quoted_string()
	}
	p.SetState(182)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SyslParserCOMMA {
		{
			p.SetState(178)
			p.Match(SyslParserCOMMA)
		}
		{
			p.SetState(179)
			p.Quoted_string()
		}

		p.SetState(184)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(185)
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
		p.SetState(187)
		p.Match(SyslParserSQ_OPEN)
	}
	{
		p.SetState(188)
		p.Array_of_strings()
	}
	{
		p.SetState(189)
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
		p.SetState(191)
		p.Match(SyslParserName)
	}
	{
		p.SetState(192)
		p.Match(SyslParserEQ)
	}
	p.SetState(196)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 4, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(193)
			p.Quoted_string()
		}

	case 2:
		{
			p.SetState(194)
			p.Array_of_strings()
		}

	case 3:
		{
			p.SetState(195)
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
		p.SetState(198)
		p.Match(SyslParserSQ_OPEN)
	}
	{
		p.SetState(199)
		p.Nvp()
	}
	p.SetState(204)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SyslParserCOMMA {
		{
			p.SetState(200)
			p.Match(SyslParserCOMMA)
		}
		{
			p.SetState(201)
			p.Nvp()
		}

		p.SetState(206)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(207)
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

	p.SetState(211)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SyslParserName:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(209)
			p.Nvp()
		}

	case SyslParserTILDE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(210)
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
		p.SetState(213)
		p.Match(SyslParserSQ_OPEN)
	}
	{
		p.SetState(214)
		p.Entry()
	}
	p.SetState(219)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SyslParserCOMMA {
		{
			p.SetState(215)
			p.Match(SyslParserCOMMA)
		}
		{
			p.SetState(216)
			p.Entry()
		}

		p.SetState(221)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(222)
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
		p.SetState(224)
		p.Match(SyslParserSET_OF)
	}
	p.SetState(225)
	_la = p.GetTokenStream().LA(1)

	if !(_la == SyslParserNativeDataTypes || _la == SyslParserName) {
		p.GetErrorHandler().RecoverInline(p)
	} else {
		p.GetErrorHandler().ReportMatch(p)
		p.Consume()
	}
	p.SetState(227)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserOPEN_PAREN {
		{
			p.SetState(226)
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
		p.SetState(229)
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
		p.SetState(231)
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
		p.SetState(233)
		p.Match(SyslParserCOLON)
	}
	{
		p.SetState(234)
		p.Match(SyslParserINDENT)
	}
	p.SetState(236)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == SyslParserPIPE {
		{
			p.SetState(235)
			p.Doc_string()
		}

		p.SetState(238)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(240)
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

	p.SetState(245)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SyslParserQSTRING:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(242)
			p.Match(SyslParserQSTRING)
		}

	case SyslParserSQ_OPEN:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(243)
			p.Array_of_strings()
		}

	case SyslParserCOLON:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(244)
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
		p.SetState(247)
		p.Match(SyslParserAT)
	}
	{
		p.SetState(248)
		p.Match(SyslParserVAR_NAME)
	}
	{
		p.SetState(249)
		p.Match(SyslParserEQ)
	}
	{
		p.SetState(250)
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
		p.SetState(252)
		p.Match(SyslParserINDENT)
	}
	p.SetState(254)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == SyslParserAT {
		{
			p.SetState(253)
			p.Annotation()
		}

		p.SetState(256)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(258)
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

	p.SetState(279)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SyslParserSET_OF:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(260)
			p.Collection_type()
		}

	case SyslParserNativeDataTypes, SyslParserName:
		p.EnterOuterAlt(localctx, 2)
		p.SetState(264)
		p.GetErrorHandler().Sync(p)
		switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 12, p.GetParserRuleContext()) {
		case 1:
			{
				p.SetState(261)
				p.Reference()
			}

		case 2:
			{
				p.SetState(262)
				p.Match(SyslParserNativeDataTypes)
			}

		case 3:
			{
				p.SetState(263)
				p.User_defined_type()
			}

		}
		p.SetState(267)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SyslParserOPEN_PAREN {
			{
				p.SetState(266)
				p.Size_spec()
			}

		}
		p.SetState(270)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SyslParserQN {
			{
				p.SetState(269)
				p.Match(SyslParserQN)
			}

		}
		p.SetState(273)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SyslParserSQ_OPEN {
			{
				p.SetState(272)
				p.Attribs_or_modifiers()
			}

		}
		p.SetState(277)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SyslParserCOLON {
			{
				p.SetState(275)
				p.Match(SyslParserCOLON)
			}
			{
				p.SetState(276)
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
		p.SetState(281)
		p.Match(SyslParserName)
	}
	{
		p.SetState(282)
		p.Match(SyslParserLESS_COLON)
	}
	{
		p.SetState(283)
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
	p.SetState(288)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SyslParserSYSL_COMMENT {
		{
			p.SetState(285)
			p.Match(SyslParserSYSL_COMMENT)
		}

		p.SetState(290)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(291)
	_la = p.GetTokenStream().LA(1)

	if !(_la == SyslParserTABLE || _la == SyslParserTYPE) {
		p.GetErrorHandler().RecoverInline(p)
	} else {
		p.GetErrorHandler().ReportMatch(p)
		p.Consume()
	}
	{
		p.SetState(292)
		p.Match(SyslParserName)
	}
	p.SetState(294)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserSQ_OPEN {
		{
			p.SetState(293)
			p.Attribs_or_modifiers()
		}

	}
	{
		p.SetState(296)
		p.Match(SyslParserCOLON)
	}
	{
		p.SetState(297)
		p.Match(SyslParserINDENT)
	}
	p.SetState(300)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == SyslParserSYSL_COMMENT || _la == SyslParserName {
		p.SetState(300)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case SyslParserSYSL_COMMENT:
			{
				p.SetState(298)
				p.Match(SyslParserSYSL_COMMENT)
			}

		case SyslParserName:
			{
				p.SetState(299)
				p.Field()
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(302)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(304)
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

func (s *Package_nameContext) AllName() []antlr.TerminalNode {
	return s.GetTokens(SyslParserName)
}

func (s *Package_nameContext) Name(i int) antlr.TerminalNode {
	return s.GetToken(SyslParserName, i)
}

func (s *Package_nameContext) DOT() antlr.TerminalNode {
	return s.GetToken(SyslParserDOT, 0)
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

	p.SetState(313)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SyslParserName:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(306)
			p.Match(SyslParserName)
		}
		p.SetState(309)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SyslParserDOT {
			{
				p.SetState(307)
				p.Match(SyslParserDOT)
			}
			{
				p.SetState(308)
				p.Match(SyslParserName)
			}

		}

	case SyslParserTEXT_LINE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(311)
			p.Match(SyslParserTEXT_LINE)
		}

	case SyslParserTEXT_NAME:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(312)
			p.Match(SyslParserTEXT_NAME)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
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
		p.SetState(315)
		p.Match(SyslParserNAME_SEP)
	}
	{
		p.SetState(316)
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
		p.SetState(318)
		p.Package_name()
	}
	p.SetState(322)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SyslParserNAME_SEP {
		{
			p.SetState(319)
			p.Sub_package()
		}

		p.SetState(324)
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
		p.SetState(325)
		p.App_name()
	}
	p.SetState(327)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserQSTRING {
		{
			p.SetState(326)
			p.Match(SyslParserQSTRING)
		}

	}
	p.SetState(330)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserSQ_OPEN {
		{
			p.SetState(329)
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
		p.SetState(332)
		p.Match(SyslParserName)
	}
	{
		p.SetState(333)
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
		p.SetState(335)
		p.Match(SyslParserCOLON)
	}
	{
		p.SetState(336)
		p.Match(SyslParserINDENT)
	}
	p.SetState(341)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == SyslParserName {
		{
			p.SetState(337)
			p.Match(SyslParserName)
		}
		p.SetState(339)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SyslParserSQ_OPEN {
			{
				p.SetState(338)
				p.Attribs_or_modifiers()
			}

		}

		p.SetState(343)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(345)
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
	p.SetState(347)
	_la = p.GetTokenStream().LA(1)

	if !(_la == SyslParserTABLE || _la == SyslParserTYPE) {
		p.GetErrorHandler().RecoverInline(p)
	} else {
		p.GetErrorHandler().ReportMatch(p)
		p.Consume()
	}
	{
		p.SetState(348)
		p.Match(SyslParserName)
	}
	p.SetState(350)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserCOLON {
		{
			p.SetState(349)
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
	p.SetState(355)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SyslParserSYSL_COMMENT {
		{
			p.SetState(352)
			p.Match(SyslParserSYSL_COMMENT)
		}

		p.SetState(357)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(358)
		p.Match(SyslParserWRAP)
	}
	{
		p.SetState(359)
		p.Model_name()
	}
	{
		p.SetState(360)
		p.Match(SyslParserINDENT)
	}
	p.SetState(362)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == SyslParserTABLE || _la == SyslParserTYPE {
		{
			p.SetState(361)
			p.Table_refs()
		}

		p.SetState(364)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(366)
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
		p.SetState(368)
		p.Match(SyslParserAT)
	}
	{
		p.SetState(369)
		p.Match(SyslParserName)
	}
	{
		p.SetState(370)
		p.Match(SyslParserEQ)
	}
	{
		p.SetState(371)
		p.Match(SyslParserQSTRING)
	}
	{
		p.SetState(372)
		p.Match(SyslParserNEWLINE)
	}

	return localctx
}

// IVar_in_curlyContext is an interface to support dynamic dispatch.
type IVar_in_curlyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsVar_in_curlyContext differentiates from other interfaces.
	IsVar_in_curlyContext()
}

type Var_in_curlyContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyVar_in_curlyContext() *Var_in_curlyContext {
	var p = new(Var_in_curlyContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_var_in_curly
	return p
}

func (*Var_in_curlyContext) IsVar_in_curlyContext() {}

func NewVar_in_curlyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Var_in_curlyContext {
	var p = new(Var_in_curlyContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_var_in_curly

	return p
}

func (s *Var_in_curlyContext) GetParser() antlr.Parser { return s.parser }

func (s *Var_in_curlyContext) CURLY_OPEN() antlr.TerminalNode {
	return s.GetToken(SyslParserCURLY_OPEN, 0)
}

func (s *Var_in_curlyContext) Name() antlr.TerminalNode {
	return s.GetToken(SyslParserName, 0)
}

func (s *Var_in_curlyContext) CURLY_CLOSE() antlr.TerminalNode {
	return s.GetToken(SyslParserCURLY_CLOSE, 0)
}

func (s *Var_in_curlyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Var_in_curlyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Var_in_curlyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterVar_in_curly(s)
	}
}

func (s *Var_in_curlyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitVar_in_curly(s)
	}
}

func (p *SyslParser) Var_in_curly() (localctx IVar_in_curlyContext) {
	localctx = NewVar_in_curlyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, SyslParserRULE_var_in_curly)

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
		p.SetState(374)
		p.Match(SyslParserCURLY_OPEN)
	}
	{
		p.SetState(375)
		p.Match(SyslParserName)
	}
	{
		p.SetState(376)
		p.Match(SyslParserCURLY_CLOSE)
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

func (s *Query_varContext) Var_in_curly() IVar_in_curlyContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVar_in_curlyContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IVar_in_curlyContext)
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
	p.EnterRule(localctx, 66, SyslParserRULE_query_var)

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
		p.SetState(378)
		p.Match(SyslParserName)
	}
	{
		p.SetState(379)
		p.Match(SyslParserEQ)
	}
	p.SetState(383)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SyslParserName:
		{
			p.SetState(380)
			p.Match(SyslParserName)
		}

	case SyslParserNativeDataTypes:
		{
			p.SetState(381)
			p.Match(SyslParserNativeDataTypes)
		}

	case SyslParserCURLY_OPEN:
		{
			p.SetState(382)
			p.Var_in_curly()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
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
	p.EnterRule(localctx, 68, SyslParserRULE_query_param)
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
		p.SetState(385)
		p.Match(SyslParserQN)
	}
	{
		p.SetState(386)
		p.Query_var()
	}
	p.SetState(391)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SyslParserAMP {
		{
			p.SetState(387)
			p.Match(SyslParserAMP)
		}
		{
			p.SetState(388)
			p.Query_var()
		}

		p.SetState(393)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IHttp_path_var_with_typeContext is an interface to support dynamic dispatch.
type IHttp_path_var_with_typeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsHttp_path_var_with_typeContext differentiates from other interfaces.
	IsHttp_path_var_with_typeContext()
}

type Http_path_var_with_typeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHttp_path_var_with_typeContext() *Http_path_var_with_typeContext {
	var p = new(Http_path_var_with_typeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_http_path_var_with_type
	return p
}

func (*Http_path_var_with_typeContext) IsHttp_path_var_with_typeContext() {}

func NewHttp_path_var_with_typeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Http_path_var_with_typeContext {
	var p = new(Http_path_var_with_typeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_http_path_var_with_type

	return p
}

func (s *Http_path_var_with_typeContext) GetParser() antlr.Parser { return s.parser }

func (s *Http_path_var_with_typeContext) CURLY_OPEN() antlr.TerminalNode {
	return s.GetToken(SyslParserCURLY_OPEN, 0)
}

func (s *Http_path_var_with_typeContext) AllName() []antlr.TerminalNode {
	return s.GetTokens(SyslParserName)
}

func (s *Http_path_var_with_typeContext) Name(i int) antlr.TerminalNode {
	return s.GetToken(SyslParserName, i)
}

func (s *Http_path_var_with_typeContext) LESS_COLON() antlr.TerminalNode {
	return s.GetToken(SyslParserLESS_COLON, 0)
}

func (s *Http_path_var_with_typeContext) CURLY_CLOSE() antlr.TerminalNode {
	return s.GetToken(SyslParserCURLY_CLOSE, 0)
}

func (s *Http_path_var_with_typeContext) NativeDataTypes() antlr.TerminalNode {
	return s.GetToken(SyslParserNativeDataTypes, 0)
}

func (s *Http_path_var_with_typeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Http_path_var_with_typeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Http_path_var_with_typeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterHttp_path_var_with_type(s)
	}
}

func (s *Http_path_var_with_typeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitHttp_path_var_with_type(s)
	}
}

func (p *SyslParser) Http_path_var_with_type() (localctx IHttp_path_var_with_typeContext) {
	localctx = NewHttp_path_var_with_typeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 70, SyslParserRULE_http_path_var_with_type)
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
		p.SetState(394)
		p.Match(SyslParserCURLY_OPEN)
	}
	{
		p.SetState(395)
		p.Match(SyslParserName)
	}
	{
		p.SetState(396)
		p.Match(SyslParserLESS_COLON)
	}
	p.SetState(397)
	_la = p.GetTokenStream().LA(1)

	if !(_la == SyslParserNativeDataTypes || _la == SyslParserName) {
		p.GetErrorHandler().RecoverInline(p)
	} else {
		p.GetErrorHandler().ReportMatch(p)
		p.Consume()
	}
	{
		p.SetState(398)
		p.Match(SyslParserCURLY_CLOSE)
	}

	return localctx
}

// IHttp_path_staticContext is an interface to support dynamic dispatch.
type IHttp_path_staticContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsHttp_path_staticContext differentiates from other interfaces.
	IsHttp_path_staticContext()
}

type Http_path_staticContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHttp_path_staticContext() *Http_path_staticContext {
	var p = new(Http_path_staticContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_http_path_static
	return p
}

func (*Http_path_staticContext) IsHttp_path_staticContext() {}

func NewHttp_path_staticContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Http_path_staticContext {
	var p = new(Http_path_staticContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_http_path_static

	return p
}

func (s *Http_path_staticContext) GetParser() antlr.Parser { return s.parser }

func (s *Http_path_staticContext) Name() antlr.TerminalNode {
	return s.GetToken(SyslParserName, 0)
}

func (s *Http_path_staticContext) TEXT_LINE() antlr.TerminalNode {
	return s.GetToken(SyslParserTEXT_LINE, 0)
}

func (s *Http_path_staticContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Http_path_staticContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Http_path_staticContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterHttp_path_static(s)
	}
}

func (s *Http_path_staticContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitHttp_path_static(s)
	}
}

func (p *SyslParser) Http_path_static() (localctx IHttp_path_staticContext) {
	localctx = NewHttp_path_staticContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 72, SyslParserRULE_http_path_static)
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
	p.SetState(400)
	_la = p.GetTokenStream().LA(1)

	if !(_la == SyslParserTEXT_LINE || _la == SyslParserName) {
		p.GetErrorHandler().RecoverInline(p)
	} else {
		p.GetErrorHandler().ReportMatch(p)
		p.Consume()
	}

	return localctx
}

// IHttp_path_suffixContext is an interface to support dynamic dispatch.
type IHttp_path_suffixContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsHttp_path_suffixContext differentiates from other interfaces.
	IsHttp_path_suffixContext()
}

type Http_path_suffixContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHttp_path_suffixContext() *Http_path_suffixContext {
	var p = new(Http_path_suffixContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_http_path_suffix
	return p
}

func (*Http_path_suffixContext) IsHttp_path_suffixContext() {}

func NewHttp_path_suffixContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Http_path_suffixContext {
	var p = new(Http_path_suffixContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_http_path_suffix

	return p
}

func (s *Http_path_suffixContext) GetParser() antlr.Parser { return s.parser }

func (s *Http_path_suffixContext) FORWARD_SLASH() antlr.TerminalNode {
	return s.GetToken(SyslParserFORWARD_SLASH, 0)
}

func (s *Http_path_suffixContext) Http_path_static() IHttp_path_staticContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHttp_path_staticContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IHttp_path_staticContext)
}

func (s *Http_path_suffixContext) Http_path_var_with_type() IHttp_path_var_with_typeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHttp_path_var_with_typeContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IHttp_path_var_with_typeContext)
}

func (s *Http_path_suffixContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Http_path_suffixContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Http_path_suffixContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterHttp_path_suffix(s)
	}
}

func (s *Http_path_suffixContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitHttp_path_suffix(s)
	}
}

func (p *SyslParser) Http_path_suffix() (localctx IHttp_path_suffixContext) {
	localctx = NewHttp_path_suffixContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 74, SyslParserRULE_http_path_suffix)

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
		p.SetState(402)
		p.Match(SyslParserFORWARD_SLASH)
	}
	p.SetState(405)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SyslParserTEXT_LINE, SyslParserName:
		{
			p.SetState(403)
			p.Http_path_static()
		}

	case SyslParserCURLY_OPEN:
		{
			p.SetState(404)
			p.Http_path_var_with_type()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
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

func (s *Http_pathContext) FORWARD_SLASH() antlr.TerminalNode {
	return s.GetToken(SyslParserFORWARD_SLASH, 0)
}

func (s *Http_pathContext) AllHttp_path_suffix() []IHttp_path_suffixContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IHttp_path_suffixContext)(nil)).Elem())
	var tst = make([]IHttp_path_suffixContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IHttp_path_suffixContext)
		}
	}

	return tst
}

func (s *Http_pathContext) Http_path_suffix(i int) IHttp_path_suffixContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHttp_path_suffixContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IHttp_path_suffixContext)
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
	p.EnterRule(localctx, 76, SyslParserRULE_http_path)
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

	p.SetState(413)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 36, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(407)
			p.Match(SyslParserFORWARD_SLASH)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		p.SetState(409)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == SyslParserFORWARD_SLASH {
			{
				p.SetState(408)
				p.Http_path_suffix()
			}

			p.SetState(411)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
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

func (s *Endpoint_nameContext) TEXT_LINE() antlr.TerminalNode {
	return s.GetToken(SyslParserTEXT_LINE, 0)
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
	p.EnterRule(localctx, 78, SyslParserRULE_endpoint_name)
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
	p.SetState(424)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SyslParserDOT, SyslParserName:
		p.SetState(416)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SyslParserDOT {
			{
				p.SetState(415)
				p.Match(SyslParserDOT)
			}

		}
		p.SetState(419)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == SyslParserName {
			{
				p.SetState(418)
				p.Match(SyslParserName)
			}

			p.SetState(421)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case SyslParserTEXT_LINE:
		{
			p.SetState(423)
			p.Match(SyslParserTEXT_LINE)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
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
	p.EnterRule(localctx, 80, SyslParserRULE_ret_stmt)

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
		p.SetState(426)
		p.Match(SyslParserRETURN)
	}
	{
		p.SetState(427)
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
	p.EnterRule(localctx, 82, SyslParserRULE_target)

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
	p.SetState(431)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SyslParserDOT:
		{
			p.SetState(429)
			p.Match(SyslParserDOT)
		}

	case SyslParserTEXT_LINE, SyslParserName, SyslParserTEXT_NAME:
		{
			p.SetState(430)
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
	p.EnterRule(localctx, 84, SyslParserRULE_target_endpoint)
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
	p.SetState(433)
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

func (s *Call_stmtContext) ARROW_LEFT() antlr.TerminalNode {
	return s.GetToken(SyslParserARROW_LEFT, 0)
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
	p.EnterRule(localctx, 86, SyslParserRULE_call_stmt)

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
		p.SetState(435)
		p.Target()
	}

	{
		p.SetState(436)
		p.Match(SyslParserARROW_LEFT)
	}
	{
		p.SetState(437)
		p.Target_endpoint()
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

func (s *If_stmtContext) AllStatements() []IStatementsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IStatementsContext)(nil)).Elem())
	var tst = make([]IStatementsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IStatementsContext)
		}
	}

	return tst
}

func (s *If_stmtContext) Statements(i int) IStatementsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStatementsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IStatementsContext)
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
	p.EnterRule(localctx, 88, SyslParserRULE_if_stmt)
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
		p.SetState(439)
		p.Match(SyslParserIF)
	}
	{
		p.SetState(440)
		p.Match(SyslParserTEXT_NAME)
	}
	{
		p.SetState(441)
		p.Match(SyslParserCOLON)
	}
	{
		p.SetState(442)
		p.Match(SyslParserINDENT)
	}
	p.SetState(446)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<SyslParserRETURN)|(1<<SyslParserIF)|(1<<SyslParserFOR)|(1<<SyslParserWHATEVER)|(1<<SyslParserONE_OF))) != 0) || (((_la-33)&-(0x1f+1)) == 0 && ((1<<uint((_la-33)))&((1<<(SyslParserDOT-33))|(1<<(SyslParserAT-33))|(1<<(SyslParserPIPE-33))|(1<<(SyslParserQSTRING-33))|(1<<(SyslParserSYSL_COMMENT-33))|(1<<(SyslParserTEXT_LINE-33))|(1<<(SyslParserName-33))|(1<<(SyslParserTEXT_NAME-33)))) != 0) {
		{
			p.SetState(443)
			p.Statements()
		}

		p.SetState(448)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(449)
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

func (s *If_elseContext) AllStatements() []IStatementsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IStatementsContext)(nil)).Elem())
	var tst = make([]IStatementsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IStatementsContext)
		}
	}

	return tst
}

func (s *If_elseContext) Statements(i int) IStatementsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStatementsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IStatementsContext)
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
	p.EnterRule(localctx, 90, SyslParserRULE_if_else)
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
		p.SetState(451)
		p.If_stmt()
	}
	p.SetState(462)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserELSE {
		{
			p.SetState(452)
			p.Match(SyslParserELSE)
		}
		{
			p.SetState(453)
			p.Match(SyslParserCOLON)
		}
		{
			p.SetState(454)
			p.Match(SyslParserINDENT)
		}
		p.SetState(458)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<SyslParserRETURN)|(1<<SyslParserIF)|(1<<SyslParserFOR)|(1<<SyslParserWHATEVER)|(1<<SyslParserONE_OF))) != 0) || (((_la-33)&-(0x1f+1)) == 0 && ((1<<uint((_la-33)))&((1<<(SyslParserDOT-33))|(1<<(SyslParserAT-33))|(1<<(SyslParserPIPE-33))|(1<<(SyslParserQSTRING-33))|(1<<(SyslParserSYSL_COMMENT-33))|(1<<(SyslParserTEXT_LINE-33))|(1<<(SyslParserName-33))|(1<<(SyslParserTEXT_NAME-33)))) != 0) {
			{
				p.SetState(455)
				p.Statements()
			}

			p.SetState(460)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(461)
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
	p.EnterRule(localctx, 92, SyslParserRULE_for_cond)

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
		p.SetState(464)
		p.Match(SyslParserTEXT_NAME)
	}
	{
		p.SetState(465)
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

func (s *For_stmtContext) AllStatements() []IStatementsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IStatementsContext)(nil)).Elem())
	var tst = make([]IStatementsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IStatementsContext)
		}
	}

	return tst
}

func (s *For_stmtContext) Statements(i int) IStatementsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStatementsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IStatementsContext)
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
	p.EnterRule(localctx, 94, SyslParserRULE_for_stmt)
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
		p.SetState(467)
		p.Match(SyslParserFOR)
	}
	{
		p.SetState(468)
		p.For_cond()
	}
	{
		p.SetState(469)
		p.Match(SyslParserINDENT)
	}
	p.SetState(473)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<SyslParserRETURN)|(1<<SyslParserIF)|(1<<SyslParserFOR)|(1<<SyslParserWHATEVER)|(1<<SyslParserONE_OF))) != 0) || (((_la-33)&-(0x1f+1)) == 0 && ((1<<uint((_la-33)))&((1<<(SyslParserDOT-33))|(1<<(SyslParserAT-33))|(1<<(SyslParserPIPE-33))|(1<<(SyslParserQSTRING-33))|(1<<(SyslParserSYSL_COMMENT-33))|(1<<(SyslParserTEXT_LINE-33))|(1<<(SyslParserName-33))|(1<<(SyslParserTEXT_NAME-33)))) != 0) {
		{
			p.SetState(470)
			p.Statements()
		}

		p.SetState(475)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(476)
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
	p.EnterRule(localctx, 96, SyslParserRULE_http_method_comment)

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
		p.SetState(478)
		p.Match(SyslParserSYSL_COMMENT)
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
	p.EnterRule(localctx, 98, SyslParserRULE_one_of_case_label)
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
	p.SetState(483)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SyslParserName {
		{
			p.SetState(480)
			p.Match(SyslParserName)
		}

		p.SetState(485)
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

func (s *One_of_casesContext) AllStatements() []IStatementsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IStatementsContext)(nil)).Elem())
	var tst = make([]IStatementsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IStatementsContext)
		}
	}

	return tst
}

func (s *One_of_casesContext) Statements(i int) IStatementsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStatementsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IStatementsContext)
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
	p.EnterRule(localctx, 100, SyslParserRULE_one_of_cases)
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
		p.SetState(486)
		p.One_of_case_label()
	}
	{
		p.SetState(487)
		p.Match(SyslParserCOLON)
	}
	{
		p.SetState(488)
		p.Match(SyslParserINDENT)
	}
	p.SetState(490)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<SyslParserRETURN)|(1<<SyslParserIF)|(1<<SyslParserFOR)|(1<<SyslParserWHATEVER)|(1<<SyslParserONE_OF))) != 0) || (((_la-33)&-(0x1f+1)) == 0 && ((1<<uint((_la-33)))&((1<<(SyslParserDOT-33))|(1<<(SyslParserAT-33))|(1<<(SyslParserPIPE-33))|(1<<(SyslParserQSTRING-33))|(1<<(SyslParserSYSL_COMMENT-33))|(1<<(SyslParserTEXT_LINE-33))|(1<<(SyslParserName-33))|(1<<(SyslParserTEXT_NAME-33)))) != 0) {
		{
			p.SetState(489)
			p.Statements()
		}

		p.SetState(492)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(494)
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
	p.EnterRule(localctx, 102, SyslParserRULE_one_of_stmt)
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
		p.SetState(496)
		p.Match(SyslParserONE_OF)
	}
	{
		p.SetState(497)
		p.Match(SyslParserCOLON)
	}
	{
		p.SetState(498)
		p.Match(SyslParserINDENT)
	}
	p.SetState(500)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == SyslParserCOLON || _la == SyslParserName {
		{
			p.SetState(499)
			p.One_of_cases()
		}

		p.SetState(502)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(504)
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

func (s *Text_stmtContext) Doc_string() IDoc_stringContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDoc_stringContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDoc_stringContext)
}

func (s *Text_stmtContext) TEXT_LINE() antlr.TerminalNode {
	return s.GetToken(SyslParserTEXT_LINE, 0)
}

func (s *Text_stmtContext) WHATEVER() antlr.TerminalNode {
	return s.GetToken(SyslParserWHATEVER, 0)
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
	p.EnterRule(localctx, 104, SyslParserRULE_text_stmt)

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

	p.SetState(509)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SyslParserPIPE:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(506)
			p.Doc_string()
		}

	case SyslParserTEXT_LINE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(507)
			p.Match(SyslParserTEXT_LINE)
		}

	case SyslParserWHATEVER:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(508)
			p.Match(SyslParserWHATEVER)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IStatementsContext is an interface to support dynamic dispatch.
type IStatementsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStatementsContext differentiates from other interfaces.
	IsStatementsContext()
}

type StatementsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStatementsContext() *StatementsContext {
	var p = new(StatementsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_statements
	return p
}

func (*StatementsContext) IsStatementsContext() {}

func NewStatementsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatementsContext {
	var p = new(StatementsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_statements

	return p
}

func (s *StatementsContext) GetParser() antlr.Parser { return s.parser }

func (s *StatementsContext) If_else() IIf_elseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIf_elseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIf_elseContext)
}

func (s *StatementsContext) For_stmt() IFor_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFor_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFor_stmtContext)
}

func (s *StatementsContext) Ret_stmt() IRet_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IRet_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IRet_stmtContext)
}

func (s *StatementsContext) Call_stmt() ICall_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICall_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICall_stmtContext)
}

func (s *StatementsContext) One_of_stmt() IOne_of_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IOne_of_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IOne_of_stmtContext)
}

func (s *StatementsContext) Http_method_comment() IHttp_method_commentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHttp_method_commentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IHttp_method_commentContext)
}

func (s *StatementsContext) QSTRING() antlr.TerminalNode {
	return s.GetToken(SyslParserQSTRING, 0)
}

func (s *StatementsContext) Text_stmt() IText_stmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IText_stmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IText_stmtContext)
}

func (s *StatementsContext) Annotation() IAnnotationContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAnnotationContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAnnotationContext)
}

func (s *StatementsContext) Attribs_or_modifiers() IAttribs_or_modifiersContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAttribs_or_modifiersContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAttribs_or_modifiersContext)
}

func (s *StatementsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StatementsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StatementsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterStatements(s)
	}
}

func (s *StatementsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitStatements(s)
	}
}

func (p *SyslParser) Statements() (localctx IStatementsContext) {
	localctx = NewStatementsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 106, SyslParserRULE_statements)
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
	p.SetState(520)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 49, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(511)
			p.If_else()
		}

	case 2:
		{
			p.SetState(512)
			p.For_stmt()
		}

	case 3:
		{
			p.SetState(513)
			p.Ret_stmt()
		}

	case 4:
		{
			p.SetState(514)
			p.Call_stmt()
		}

	case 5:
		{
			p.SetState(515)
			p.One_of_stmt()
		}

	case 6:
		{
			p.SetState(516)
			p.Http_method_comment()
		}

	case 7:
		{
			p.SetState(517)
			p.Match(SyslParserQSTRING)
		}

	case 8:
		{
			p.SetState(518)
			p.Text_stmt()
		}

	case 9:
		{
			p.SetState(519)
			p.Annotation()
		}

	}
	p.SetState(523)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserSQ_OPEN {
		{
			p.SetState(522)
			p.Attribs_or_modifiers()
		}

	}

	return localctx
}

// IMethod_defContext is an interface to support dynamic dispatch.
type IMethod_defContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMethod_defContext differentiates from other interfaces.
	IsMethod_defContext()
}

type Method_defContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMethod_defContext() *Method_defContext {
	var p = new(Method_defContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_method_def
	return p
}

func (*Method_defContext) IsMethod_defContext() {}

func NewMethod_defContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Method_defContext {
	var p = new(Method_defContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_method_def

	return p
}

func (s *Method_defContext) GetParser() antlr.Parser { return s.parser }

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

func (s *Method_defContext) AllStatements() []IStatementsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IStatementsContext)(nil)).Elem())
	var tst = make([]IStatementsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IStatementsContext)
		}
	}

	return tst
}

func (s *Method_defContext) Statements(i int) IStatementsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStatementsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IStatementsContext)
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

func (p *SyslParser) Method_def() (localctx IMethod_defContext) {
	localctx = NewMethod_defContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 108, SyslParserRULE_method_def)
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
		p.SetState(525)
		p.Match(SyslParserHTTP_VERBS)
	}
	p.SetState(527)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserQN {
		{
			p.SetState(526)
			p.Query_param()
		}

	}
	p.SetState(530)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserSQ_OPEN {
		{
			p.SetState(529)
			p.Attributes()
		}

	}
	{
		p.SetState(532)
		p.Match(SyslParserCOLON)
	}
	{
		p.SetState(533)
		p.Match(SyslParserINDENT)
	}
	p.SetState(535)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<SyslParserRETURN)|(1<<SyslParserIF)|(1<<SyslParserFOR)|(1<<SyslParserWHATEVER)|(1<<SyslParserONE_OF))) != 0) || (((_la-33)&-(0x1f+1)) == 0 && ((1<<uint((_la-33)))&((1<<(SyslParserDOT-33))|(1<<(SyslParserAT-33))|(1<<(SyslParserPIPE-33))|(1<<(SyslParserQSTRING-33))|(1<<(SyslParserSYSL_COMMENT-33))|(1<<(SyslParserTEXT_LINE-33))|(1<<(SyslParserName-33))|(1<<(SyslParserTEXT_NAME-33)))) != 0) {
		{
			p.SetState(534)
			p.Statements()
		}

		p.SetState(537)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(539)
		p.Match(SyslParserDEDENT)
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
		p.SetState(541)
		p.Match(SyslParserWHATEVER)
	}

	return localctx
}

// ISimple_endpointContext is an interface to support dynamic dispatch.
type ISimple_endpointContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSimple_endpointContext differentiates from other interfaces.
	IsSimple_endpointContext()
}

type Simple_endpointContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySimple_endpointContext() *Simple_endpointContext {
	var p = new(Simple_endpointContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_simple_endpoint
	return p
}

func (*Simple_endpointContext) IsSimple_endpointContext() {}

func NewSimple_endpointContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Simple_endpointContext {
	var p = new(Simple_endpointContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_simple_endpoint

	return p
}

func (s *Simple_endpointContext) GetParser() antlr.Parser { return s.parser }

func (s *Simple_endpointContext) WHATEVER() antlr.TerminalNode {
	return s.GetToken(SyslParserWHATEVER, 0)
}

func (s *Simple_endpointContext) Endpoint_name() IEndpoint_nameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEndpoint_nameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IEndpoint_nameContext)
}

func (s *Simple_endpointContext) COLON() antlr.TerminalNode {
	return s.GetToken(SyslParserCOLON, 0)
}

func (s *Simple_endpointContext) Shortcut() IShortcutContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShortcutContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShortcutContext)
}

func (s *Simple_endpointContext) QSTRING() antlr.TerminalNode {
	return s.GetToken(SyslParserQSTRING, 0)
}

func (s *Simple_endpointContext) Attribs_or_modifiers() IAttribs_or_modifiersContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAttribs_or_modifiersContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAttribs_or_modifiersContext)
}

func (s *Simple_endpointContext) INDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserINDENT, 0)
}

func (s *Simple_endpointContext) DEDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserDEDENT, 0)
}

func (s *Simple_endpointContext) AllStatements() []IStatementsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IStatementsContext)(nil)).Elem())
	var tst = make([]IStatementsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IStatementsContext)
		}
	}

	return tst
}

func (s *Simple_endpointContext) Statements(i int) IStatementsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStatementsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IStatementsContext)
}

func (s *Simple_endpointContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Simple_endpointContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Simple_endpointContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterSimple_endpoint(s)
	}
}

func (s *Simple_endpointContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitSimple_endpoint(s)
	}
}

func (p *SyslParser) Simple_endpoint() (localctx ISimple_endpointContext) {
	localctx = NewSimple_endpointContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 112, SyslParserRULE_simple_endpoint)
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

	p.SetState(563)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SyslParserWHATEVER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(543)
			p.Match(SyslParserWHATEVER)
		}

	case SyslParserDOT, SyslParserTEXT_LINE, SyslParserName:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(544)
			p.Endpoint_name()
		}
		p.SetState(546)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SyslParserQSTRING {
			{
				p.SetState(545)
				p.Match(SyslParserQSTRING)
			}

		}
		p.SetState(549)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SyslParserSQ_OPEN {
			{
				p.SetState(548)
				p.Attribs_or_modifiers()
			}

		}
		{
			p.SetState(551)
			p.Match(SyslParserCOLON)
		}
		p.SetState(561)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case SyslParserWHATEVER:
			{
				p.SetState(552)
				p.Shortcut()
			}

		case SyslParserINDENT:
			{
				p.SetState(553)
				p.Match(SyslParserINDENT)
			}
			p.SetState(555)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			for ok := true; ok; ok = (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<SyslParserRETURN)|(1<<SyslParserIF)|(1<<SyslParserFOR)|(1<<SyslParserWHATEVER)|(1<<SyslParserONE_OF))) != 0) || (((_la-33)&-(0x1f+1)) == 0 && ((1<<uint((_la-33)))&((1<<(SyslParserDOT-33))|(1<<(SyslParserAT-33))|(1<<(SyslParserPIPE-33))|(1<<(SyslParserQSTRING-33))|(1<<(SyslParserSYSL_COMMENT-33))|(1<<(SyslParserTEXT_LINE-33))|(1<<(SyslParserName-33))|(1<<(SyslParserTEXT_NAME-33)))) != 0) {
				{
					p.SetState(554)
					p.Statements()
				}

				p.SetState(557)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(559)
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

// IRest_endpointContext is an interface to support dynamic dispatch.
type IRest_endpointContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsRest_endpointContext differentiates from other interfaces.
	IsRest_endpointContext()
}

type Rest_endpointContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRest_endpointContext() *Rest_endpointContext {
	var p = new(Rest_endpointContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_rest_endpoint
	return p
}

func (*Rest_endpointContext) IsRest_endpointContext() {}

func NewRest_endpointContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Rest_endpointContext {
	var p = new(Rest_endpointContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_rest_endpoint

	return p
}

func (s *Rest_endpointContext) GetParser() antlr.Parser { return s.parser }

func (s *Rest_endpointContext) Http_path() IHttp_pathContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHttp_pathContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IHttp_pathContext)
}

func (s *Rest_endpointContext) COLON() antlr.TerminalNode {
	return s.GetToken(SyslParserCOLON, 0)
}

func (s *Rest_endpointContext) INDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserINDENT, 0)
}

func (s *Rest_endpointContext) DEDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserDEDENT, 0)
}

func (s *Rest_endpointContext) Attribs_or_modifiers() IAttribs_or_modifiersContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAttribs_or_modifiersContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAttribs_or_modifiersContext)
}

func (s *Rest_endpointContext) AllMethod_def() []IMethod_defContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IMethod_defContext)(nil)).Elem())
	var tst = make([]IMethod_defContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IMethod_defContext)
		}
	}

	return tst
}

func (s *Rest_endpointContext) Method_def(i int) IMethod_defContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMethod_defContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IMethod_defContext)
}

func (s *Rest_endpointContext) AllRest_endpoint() []IRest_endpointContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IRest_endpointContext)(nil)).Elem())
	var tst = make([]IRest_endpointContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IRest_endpointContext)
		}
	}

	return tst
}

func (s *Rest_endpointContext) Rest_endpoint(i int) IRest_endpointContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IRest_endpointContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IRest_endpointContext)
}

func (s *Rest_endpointContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Rest_endpointContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Rest_endpointContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterRest_endpoint(s)
	}
}

func (s *Rest_endpointContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitRest_endpoint(s)
	}
}

func (p *SyslParser) Rest_endpoint() (localctx IRest_endpointContext) {
	localctx = NewRest_endpointContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 114, SyslParserRULE_rest_endpoint)
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
		p.SetState(565)
		p.Http_path()
	}
	p.SetState(567)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserSQ_OPEN {
		{
			p.SetState(566)
			p.Attribs_or_modifiers()
		}

	}
	{
		p.SetState(569)
		p.Match(SyslParserCOLON)
	}

	{
		p.SetState(570)
		p.Match(SyslParserINDENT)
	}
	p.SetState(573)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == SyslParserHTTP_VERBS || _la == SyslParserFORWARD_SLASH {
		p.SetState(573)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case SyslParserHTTP_VERBS:
			{
				p.SetState(571)
				p.Method_def()
			}

		case SyslParserFORWARD_SLASH:
			{
				p.SetState(572)
				p.Rest_endpoint()
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(575)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(577)
		p.Match(SyslParserDEDENT)
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
	p.EnterRule(localctx, 116, SyslParserRULE_collector_stmt)

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

	p.SetState(582)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SyslParserDOT, SyslParserTEXT_LINE, SyslParserName, SyslParserTEXT_NAME:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(579)
			p.Call_stmt()
		}

	case SyslParserHTTP_VERBS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(580)
			p.Match(SyslParserHTTP_VERBS)
		}
		{
			p.SetState(581)
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
	p.EnterRule(localctx, 118, SyslParserRULE_collector_stmts)

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
		p.SetState(584)
		p.Collector_stmt()
	}
	{
		p.SetState(585)
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
	p.EnterRule(localctx, 120, SyslParserRULE_collector)
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
		p.Match(SyslParserCOLLECTOR)
	}
	{
		p.SetState(588)
		p.Match(SyslParserCOLON)
	}
	p.SetState(598)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SyslParserWHATEVER:
		{
			p.SetState(589)
			p.Match(SyslParserWHATEVER)
		}

	case SyslParserINDENT:
		{
			p.SetState(590)
			p.Match(SyslParserINDENT)
		}
		p.SetState(592)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == SyslParserHTTP_VERBS || (((_la-33)&-(0x1f+1)) == 0 && ((1<<uint((_la-33)))&((1<<(SyslParserDOT-33))|(1<<(SyslParserTEXT_LINE-33))|(1<<(SyslParserName-33))|(1<<(SyslParserTEXT_NAME-33)))) != 0) {
			{
				p.SetState(591)
				p.Collector_stmts()
			}

			p.SetState(594)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(596)
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

func (s *EventContext) AllStatements() []IStatementsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IStatementsContext)(nil)).Elem())
	var tst = make([]IStatementsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IStatementsContext)
		}
	}

	return tst
}

func (s *EventContext) Statements(i int) IStatementsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStatementsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IStatementsContext)
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
	p.EnterRule(localctx, 122, SyslParserRULE_event)
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
		p.SetState(600)
		p.Match(SyslParserDISTANCE)
	}
	{
		p.SetState(601)
		p.Match(SyslParserEVENT_NAME)
	}
	p.SetState(603)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserSQ_OPEN {
		{
			p.SetState(602)
			p.Attribs_or_modifiers()
		}

	}
	{
		p.SetState(605)
		p.Match(SyslParserCOLON)
	}
	p.SetState(615)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SyslParserWHATEVER:
		{
			p.SetState(606)
			p.Match(SyslParserWHATEVER)
		}

	case SyslParserINDENT:
		{
			p.SetState(607)
			p.Match(SyslParserINDENT)
		}
		p.SetState(609)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<SyslParserRETURN)|(1<<SyslParserIF)|(1<<SyslParserFOR)|(1<<SyslParserWHATEVER)|(1<<SyslParserONE_OF))) != 0) || (((_la-33)&-(0x1f+1)) == 0 && ((1<<uint((_la-33)))&((1<<(SyslParserDOT-33))|(1<<(SyslParserAT-33))|(1<<(SyslParserPIPE-33))|(1<<(SyslParserQSTRING-33))|(1<<(SyslParserSYSL_COMMENT-33))|(1<<(SyslParserTEXT_LINE-33))|(1<<(SyslParserName-33))|(1<<(SyslParserTEXT_NAME-33)))) != 0) {
			{
				p.SetState(608)
				p.Statements()
			}

			p.SetState(611)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(613)
			p.Match(SyslParserDEDENT)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// ISubscribeContext is an interface to support dynamic dispatch.
type ISubscribeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSubscribeContext differentiates from other interfaces.
	IsSubscribeContext()
}

type SubscribeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySubscribeContext() *SubscribeContext {
	var p = new(SubscribeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SyslParserRULE_subscribe
	return p
}

func (*SubscribeContext) IsSubscribeContext() {}

func NewSubscribeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SubscribeContext {
	var p = new(SubscribeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SyslParserRULE_subscribe

	return p
}

func (s *SubscribeContext) GetParser() antlr.Parser { return s.parser }

func (s *SubscribeContext) Name() antlr.TerminalNode {
	return s.GetToken(SyslParserName, 0)
}

func (s *SubscribeContext) ARROW_RIGHT() antlr.TerminalNode {
	return s.GetToken(SyslParserARROW_RIGHT, 0)
}

func (s *SubscribeContext) COLON() antlr.TerminalNode {
	return s.GetToken(SyslParserCOLON, 0)
}

func (s *SubscribeContext) WHATEVER() antlr.TerminalNode {
	return s.GetToken(SyslParserWHATEVER, 0)
}

func (s *SubscribeContext) INDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserINDENT, 0)
}

func (s *SubscribeContext) DEDENT() antlr.TerminalNode {
	return s.GetToken(SyslParserDEDENT, 0)
}

func (s *SubscribeContext) Attribs_or_modifiers() IAttribs_or_modifiersContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAttribs_or_modifiersContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAttribs_or_modifiersContext)
}

func (s *SubscribeContext) AllStatements() []IStatementsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IStatementsContext)(nil)).Elem())
	var tst = make([]IStatementsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IStatementsContext)
		}
	}

	return tst
}

func (s *SubscribeContext) Statements(i int) IStatementsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStatementsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IStatementsContext)
}

func (s *SubscribeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SubscribeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SubscribeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.EnterSubscribe(s)
	}
}

func (s *SubscribeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SyslParserListener); ok {
		listenerT.ExitSubscribe(s)
	}
}

func (p *SyslParser) Subscribe() (localctx ISubscribeContext) {
	localctx = NewSubscribeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 124, SyslParserRULE_subscribe)
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
		p.SetState(617)
		p.Match(SyslParserName)
	}
	{
		p.SetState(618)
		p.Match(SyslParserARROW_RIGHT)
	}
	p.SetState(620)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserSQ_OPEN {
		{
			p.SetState(619)
			p.Attribs_or_modifiers()
		}

	}
	{
		p.SetState(622)
		p.Match(SyslParserCOLON)
	}
	p.SetState(632)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SyslParserWHATEVER:
		{
			p.SetState(623)
			p.Match(SyslParserWHATEVER)
		}

	case SyslParserINDENT:
		{
			p.SetState(624)
			p.Match(SyslParserINDENT)
		}
		p.SetState(626)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<SyslParserRETURN)|(1<<SyslParserIF)|(1<<SyslParserFOR)|(1<<SyslParserWHATEVER)|(1<<SyslParserONE_OF))) != 0) || (((_la-33)&-(0x1f+1)) == 0 && ((1<<uint((_la-33)))&((1<<(SyslParserDOT-33))|(1<<(SyslParserAT-33))|(1<<(SyslParserPIPE-33))|(1<<(SyslParserQSTRING-33))|(1<<(SyslParserSYSL_COMMENT-33))|(1<<(SyslParserTEXT_LINE-33))|(1<<(SyslParserName-33))|(1<<(SyslParserTEXT_NAME-33)))) != 0) {
			{
				p.SetState(625)
				p.Statements()
			}

			p.SetState(628)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(630)
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

func (s *App_declContext) AllRest_endpoint() []IRest_endpointContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IRest_endpointContext)(nil)).Elem())
	var tst = make([]IRest_endpointContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IRest_endpointContext)
		}
	}

	return tst
}

func (s *App_declContext) Rest_endpoint(i int) IRest_endpointContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IRest_endpointContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IRest_endpointContext)
}

func (s *App_declContext) AllSimple_endpoint() []ISimple_endpointContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ISimple_endpointContext)(nil)).Elem())
	var tst = make([]ISimple_endpointContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ISimple_endpointContext)
		}
	}

	return tst
}

func (s *App_declContext) Simple_endpoint(i int) ISimple_endpointContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISimple_endpointContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ISimple_endpointContext)
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

func (s *App_declContext) AllSubscribe() []ISubscribeContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ISubscribeContext)(nil)).Elem())
	var tst = make([]ISubscribeContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ISubscribeContext)
		}
	}

	return tst
}

func (s *App_declContext) Subscribe(i int) ISubscribeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISubscribeContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ISubscribeContext)
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
	p.EnterRule(localctx, 126, SyslParserRULE_app_decl)
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
		p.SetState(634)
		p.Match(SyslParserINDENT)
	}
	p.SetState(644)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<SyslParserWRAP)|(1<<SyslParserTABLE)|(1<<SyslParserTYPE)|(1<<SyslParserWHATEVER)|(1<<SyslParserDISTANCE)|(1<<SyslParserCOLLECTOR)|(1<<SyslParserFORWARD_SLASH))) != 0) || (((_la-33)&-(0x1f+1)) == 0 && ((1<<uint((_la-33)))&((1<<(SyslParserDOT-33))|(1<<(SyslParserAT-33))|(1<<(SyslParserSYSL_COMMENT-33))|(1<<(SyslParserTEXT_LINE-33))|(1<<(SyslParserName-33)))) != 0) {
		p.SetState(644)
		p.GetErrorHandler().Sync(p)
		switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 71, p.GetParserRuleContext()) {
		case 1:
			{
				p.SetState(635)
				p.Table()
			}

		case 2:
			{
				p.SetState(636)
				p.Facade()
			}

		case 3:
			{
				p.SetState(637)
				p.Match(SyslParserSYSL_COMMENT)
			}

		case 4:
			{
				p.SetState(638)
				p.Rest_endpoint()
			}

		case 5:
			{
				p.SetState(639)
				p.Simple_endpoint()
			}

		case 6:
			{
				p.SetState(640)
				p.Collector()
			}

		case 7:
			{
				p.SetState(641)
				p.Event()
			}

		case 8:
			{
				p.SetState(642)
				p.Subscribe()
			}

		case 9:
			{
				p.SetState(643)
				p.Annotation()
			}

		}

		p.SetState(646)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(648)
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
	p.EnterRule(localctx, 128, SyslParserRULE_application)
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

	for _la == SyslParserSYSL_COMMENT {
		{
			p.SetState(650)
			p.Match(SyslParserSYSL_COMMENT)
		}

		p.SetState(655)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(656)
		p.Name_with_attribs()
	}
	{
		p.SetState(657)
		p.Match(SyslParserCOLON)
	}
	{
		p.SetState(658)
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
	p.EnterRule(localctx, 130, SyslParserRULE_path)
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
	p.SetState(661)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserFORWARD_SLASH {
		{
			p.SetState(660)
			p.Match(SyslParserFORWARD_SLASH)
		}

	}
	{
		p.SetState(663)
		p.Match(SyslParserName)
	}
	p.SetState(668)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SyslParserFORWARD_SLASH {
		{
			p.SetState(664)
			p.Match(SyslParserFORWARD_SLASH)
		}
		{
			p.SetState(665)
			p.Match(SyslParserName)
		}

		p.SetState(670)
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
	p.EnterRule(localctx, 132, SyslParserRULE_import_stmt)

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
		p.SetState(671)
		p.Match(SyslParserIMPORT)
	}
	p.SetState(675)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 76, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(672)
				p.Match(SyslParserSYSL_COMMENT)
			}

		}
		p.SetState(677)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 76, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 134, SyslParserRULE_imports_decl)
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
	p.SetState(679)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == SyslParserIMPORT {
		{
			p.SetState(678)
			p.Import_stmt()
		}

		p.SetState(681)
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
	p.EnterRule(localctx, 136, SyslParserRULE_sysl_file)
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
	p.SetState(684)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SyslParserIMPORT {
		{
			p.SetState(683)
			p.Imports_decl()
		}

	}
	p.SetState(687)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (((_la-53)&-(0x1f+1)) == 0 && ((1<<uint((_la-53)))&((1<<(SyslParserSYSL_COMMENT-53))|(1<<(SyslParserTEXT_LINE-53))|(1<<(SyslParserName-53))|(1<<(SyslParserTEXT_NAME-53)))) != 0) {
		{
			p.SetState(686)
			p.Application()
		}

		p.SetState(689)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(691)
		p.Match(SyslParserEOF)
	}

	return localctx
}
