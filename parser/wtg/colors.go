package wtg

import "image/color"

var Colors = map[string]color.Color{
	//Black RGBA color
	"Black": color.RGBA{R: 0, G: 0, B: 0, A: 255},
	//White RGBA color
	"White": color.RGBA{R: 255, G: 255, B: 255, A: 255},
	//Red RGBA color
	"Red": color.RGBA{R: 255, G: 0, B: 0, A: 255},
	//Lime RGBA color
	"Lime": color.RGBA{R: 0, G: 255, B: 0, A: 255},
	//Blue RGBA color
	"Blue": color.RGBA{R: 0, G: 0, B: 255, A: 255},
	//Yellow RGBA color
	"Yellow": color.RGBA{R: 255, G: 255, B: 0, A: 255},
	//Cyan RGBA color
	"Cyan": color.RGBA{R: 0, G: 255, B: 255, A: 255},
	//Aqua RGBA color
	"Aqua": color.RGBA{R: 0, G: 255, B: 255, A: 255},
	//Magenta RGBA color
	"Magenta": color.RGBA{R: 255, G: 0, B: 255, A: 255},
	//Fuchsia RGBA color
	"Fuchsia": color.RGBA{R: 255, G: 0, B: 255, A: 255},
	//Silver RGBA color
	"Silver": color.RGBA{R: 192, G: 192, B: 192, A: 255},
	//Gray RGBA color
	"Gray": color.RGBA{R: 128, G: 128, B: 128, A: 255},
	//Maroon RGBA color
	"Maroon": color.RGBA{R: 128, G: 0, B: 0, A: 255},
	//Olive RGBA color
	"Olive": color.RGBA{R: 128, G: 128, B: 0, A: 255},
	//Green RGBA color
	"Green": color.RGBA{R: 0, G: 128, B: 0, A: 255},
	//Purple RGBA color
	"Purple": color.RGBA{R: 128, G: 0, B: 128, A: 255},
	//Teal RGBA color
	"Teal": color.RGBA{R: 0, G: 128, B: 128, A: 255},
	//Navy RGBA color
	"Navy": color.RGBA{R: 0, G: 0, B: 128, A: 255},
	//DarkRed RGBA color
	"DarkRed": color.RGBA{R: 139, G: 0, B: 0, A: 255},
	//Brown RGBA color
	"Brown": color.RGBA{R: 165, G: 42, B: 42, A: 255},
	//Firebrick RGBA color
	"Firebrick": color.RGBA{R: 178, G: 34, B: 34, A: 255},
	//Crimson RGBA color
	"Crimson": color.RGBA{R: 220, G: 20, B: 60, A: 255},
	//Tomato RGBA color
	"Tomato": color.RGBA{R: 255, G: 99, B: 71, A: 255},
	//Coral RGBA color
	"Coral": color.RGBA{R: 255, G: 127, B: 80, A: 255},
	//IndianRed RGBA color
	"IndianRed": color.RGBA{R: 205, G: 92, B: 92, A: 255},
	//LightCoral RGBA color
	"LightCoral": color.RGBA{R: 240, G: 128, B: 128, A: 255},
	//DarkSalmon RGBA color
	"DarkSalmon": color.RGBA{R: 233, G: 150, B: 122, A: 255},
	//Salmon RGBA color
	"Salmon": color.RGBA{R: 250, G: 128, B: 114, A: 255},
	//LightSalmon RGBA color
	"LightSalmon": color.RGBA{R: 255, G: 160, B: 122, A: 255},
	//OrangeRed RGBA color
	"OrangeRed": color.RGBA{R: 255, G: 69, B: 0, A: 255},
	//DarkOrange RGBA color
	"DarkOrange": color.RGBA{R: 255, G: 140, B: 0, A: 255},
	//Orange RGBA color
	"Orange": color.RGBA{R: 255, G: 165, B: 0, A: 255},
	//Gold RGBA color
	"Gold": color.RGBA{R: 255, G: 215, B: 0, A: 255},
	//DarkGoldenRod RGBA color
	"DarkGoldenRod": color.RGBA{R: 184, G: 134, B: 11, A: 255},
	//GoldenRod RGBA color
	"GoldenRod": color.RGBA{R: 218, G: 165, B: 32, A: 255},
	//PaleGoldenRod RGBA color
	"PaleGoldenRod": color.RGBA{R: 238, G: 232, B: 170, A: 255},
	//DarkKhaki RGBA color
	"DarkKhaki": color.RGBA{R: 189, G: 183, B: 107, A: 255},
	//Khaki RGBA color
	"Khaki": color.RGBA{R: 240, G: 230, B: 140, A: 255},
	//YellowGreen RGBA color
	"YellowGreen": color.RGBA{R: 154, G: 205, B: 50, A: 255},
	//DarkOliveGreen RGBA color
	"DarkOliveGreen": color.RGBA{R: 85, G: 107, B: 47, A: 255},
	//OliveDrab RGBA color
	"OliveDrab": color.RGBA{R: 107, G: 142, B: 35, A: 255},
	//LawnGreen RGBA color
	"LawnGreen": color.RGBA{R: 124, G: 252, B: 0, A: 255},
	//ChartReuse RGBA color
	"ChartReuse": color.RGBA{R: 127, G: 255, B: 0, A: 255},
	//GreenYellow RGBA color
	"GreenYellow": color.RGBA{R: 173, G: 255, B: 47, A: 255},
	//DarkGreen RGBA color
	"DarkGreen": color.RGBA{R: 0, G: 100, B: 0, A: 255},
	//ForestGreen RGBA color
	"ForestGreen": color.RGBA{R: 34, G: 139, B: 34, A: 255},
	//LimeGreen RGBA color
	"LimeGreen": color.RGBA{R: 50, G: 205, B: 50, A: 255},
	//LightGreen RGBA color
	"LightGreen": color.RGBA{R: 144, G: 238, B: 144, A: 255},
	//PaleGreen RGBA color
	"PaleGreen": color.RGBA{R: 152, G: 251, B: 152, A: 255},
	//DarkSeaGreen RGBA color
	"DarkSeaGreen": color.RGBA{R: 143, G: 188, B: 143, A: 255},
	//MediumSpringGreen RGBA color
	"MediumSpringGreen": color.RGBA{R: 0, G: 250, B: 154, A: 255},
	//SpringGreen RGBA color
	"SpringGreen": color.RGBA{R: 0, G: 255, B: 127, A: 255},
	//SeaGreen RGBA color
	"SeaGreen": color.RGBA{R: 46, G: 139, B: 87, A: 255},
	//MediumAquaMarine RGBA color
	"MediumAquaMarine": color.RGBA{R: 102, G: 205, B: 170, A: 255},
	//MediumSeaGreen RGBA color
	"MediumSeaGreen": color.RGBA{R: 60, G: 179, B: 113, A: 255},
	//LightSeaGreen RGBA color
	"LightSeaGreen": color.RGBA{R: 32, G: 178, B: 170, A: 255},
	//DarkSlateGray RGBA color
	"DarkSlateGray": color.RGBA{R: 47, G: 79, B: 79, A: 255},
	//DarkCyan RGBA color
	"DarkCyan": color.RGBA{R: 0, G: 139, B: 139, A: 255},
	//LightCyan RGBA color
	"LightCyan": color.RGBA{R: 224, G: 255, B: 255, A: 255},
	//DarkTurquoise RGBA color
	"DarkTurquoise": color.RGBA{R: 0, G: 206, B: 209, A: 255},
	//Turquoise RGBA color
	"Turquoise": color.RGBA{R: 64, G: 224, B: 208, A: 255},
	//MediumTurquoise RGBA color
	"MediumTurquoise": color.RGBA{R: 72, G: 209, B: 204, A: 255},
	//PaleTurquoise RGBA color
	"PaleTurquoise": color.RGBA{R: 175, G: 238, B: 238, A: 255},
	//AquaMarine RGBA color
	"AquaMarine": color.RGBA{R: 127, G: 255, B: 212, A: 255},
	//PowderBlue RGBA color
	"PowderBlue": color.RGBA{R: 176, G: 224, B: 230, A: 255},
	//CadetBlue RGBA color
	"CadetBlue": color.RGBA{R: 95, G: 158, B: 160, A: 255},
	//SteelBlue RGBA color
	"SteelBlue": color.RGBA{R: 70, G: 130, B: 180, A: 255},
	//CornFlowerBlue RGBA color
	"CornFlowerBlue": color.RGBA{R: 100, G: 149, B: 237, A: 255},
	//DeepSkyBlue RGBA color
	"DeepSkyBlue": color.RGBA{R: 0, G: 191, B: 255, A: 255},
	//DodgerBlue RGBA color
	"DodgerBlue": color.RGBA{R: 30, G: 144, B: 255, A: 255},
	//LightBlue RGBA color
	"LightBlue": color.RGBA{R: 173, G: 216, B: 230, A: 255},
	//SkyBlue RGBA color
	"SkyBlue": color.RGBA{R: 135, G: 206, B: 235, A: 255},
	//LightSkyBlue RGBA color
	"LightSkyBlue": color.RGBA{R: 135, G: 206, B: 250, A: 255},
	//MidnightBlue RGBA color
	"MidnightBlue": color.RGBA{R: 25, G: 25, B: 112, A: 255},
	//DarkBlue RGBA color
	"DarkBlue": color.RGBA{R: 0, G: 0, B: 139, A: 255},
	//MediumBlue RGBA color
	"MediumBlue": color.RGBA{R: 0, G: 0, B: 205, A: 255},
	//RoyalBlue RGBA color
	"RoyalBlue": color.RGBA{R: 65, G: 105, B: 225, A: 255},
	//BlueViolet RGBA color
	"BlueViolet": color.RGBA{R: 138, G: 43, B: 226, A: 255},
	//Indigo RGBA color
	"Indigo": color.RGBA{R: 75, G: 0, B: 130, A: 255},
	//DarkSlateBlue RGBA color
	"DarkSlateBlue": color.RGBA{R: 72, G: 61, B: 139, A: 255},
	//SlateBlue RGBA color
	"SlateBlue": color.RGBA{R: 106, G: 90, B: 205, A: 255},
	//MediumSlateBlue RGBA color
	"MediumSlateBlue": color.RGBA{R: 123, G: 104, B: 238, A: 255},
	//MediumPurple RGBA color
	"MediumPurple": color.RGBA{R: 147, G: 112, B: 219, A: 255},
	//DarkMagenta RGBA color
	"DarkMagenta": color.RGBA{R: 139, G: 0, B: 139, A: 255},
	//DarkViolet RGBA color
	"DarkViolet": color.RGBA{R: 148, G: 0, B: 211, A: 255},
	//DarkOrchid RGBA color
	"DarkOrchid": color.RGBA{R: 153, G: 50, B: 204, A: 255},
	//MediumOrchid RGBA color
	"MediumOrchid": color.RGBA{R: 186, G: 85, B: 211, A: 255},
	//Thistle RGBA color
	"Thistle": color.RGBA{R: 216, G: 191, B: 216, A: 255},
	//Plum RGBA color
	"Plum": color.RGBA{R: 221, G: 160, B: 221, A: 255},
	//Violet RGBA color
	"Violet": color.RGBA{R: 238, G: 130, B: 238, A: 255},
	//Orchid RGBA color
	"Orchid": color.RGBA{R: 218, G: 112, B: 214, A: 255},
	//MediumVioletRed RGBA color
	"MediumVioletRed": color.RGBA{R: 199, G: 21, B: 133, A: 255},
	//PaleVioletRed RGBA color
	"PaleVioletRed": color.RGBA{R: 219, G: 112, B: 147, A: 255},
	//DeepPink RGBA color
	"DeepPink": color.RGBA{R: 255, G: 20, B: 147, A: 255},
	//HotPink RGBA color
	"HotPink": color.RGBA{R: 255, G: 105, B: 180, A: 255},
	//LightPink RGBA color
	"LightPink": color.RGBA{R: 255, G: 182, B: 193, A: 255},
	//Pink RGBA color
	"Pink": color.RGBA{R: 255, G: 192, B: 203, A: 255},
	//AntiqueWhite RGBA color
	"AntiqueWhite": color.RGBA{R: 250, G: 235, B: 215, A: 255},
	//Beige RGBA color
	"Beige": color.RGBA{R: 245, G: 245, B: 220, A: 255},
	//Bisque RGBA color
	"Bisque": color.RGBA{R: 255, G: 228, B: 196, A: 255},
	//BlanchedAlmond RGBA color
	"BlanchedAlmond": color.RGBA{R: 255, G: 235, B: 205, A: 255},
	//Wheat RGBA color
	"Wheat": color.RGBA{R: 245, G: 222, B: 179, A: 255},
	//CornSilk RGBA color
	"CornSilk": color.RGBA{R: 255, G: 248, B: 220, A: 255},
	//LemonChiffon RGBA color
	"LemonChiffon": color.RGBA{R: 255, G: 250, B: 205, A: 255},
	//LightGoldenRod RGBA color
	"LightGoldenRod": color.RGBA{R: 250, G: 250, B: 210, A: 255},
	//LightYellow RGBA color
	"LightYellow": color.RGBA{R: 255, G: 255, B: 224, A: 255},
	//SaddleBrown RGBA color
	"SaddleBrown": color.RGBA{R: 139, G: 69, B: 19, A: 255},
	//Sienna RGBA color
	"Sienna": color.RGBA{R: 160, G: 82, B: 45, A: 255},
	//Chocolate RGBA color
	"Chocolate": color.RGBA{R: 210, G: 105, B: 30, A: 255},
	//Peru RGBA color
	"Peru": color.RGBA{R: 205, G: 133, B: 63, A: 255},
	//SandyBrown RGBA color
	"SandyBrown": color.RGBA{R: 244, G: 164, B: 96, A: 255},
	//BurlyWood RGBA color
	"BurlyWood": color.RGBA{R: 222, G: 184, B: 135, A: 255},
	//Tan RGBA color
	"Tan": color.RGBA{R: 210, G: 180, B: 140, A: 255},
	//RosyBrown RGBA color
	"RosyBrown": color.RGBA{R: 188, G: 143, B: 143, A: 255},
	//Moccasin RGBA color
	"Moccasin": color.RGBA{R: 255, G: 228, B: 181, A: 255},
	//NavajoWhite RGBA color
	"NavajoWhite": color.RGBA{R: 255, G: 222, B: 173, A: 255},
	//PeachPuff RGBA color
	"PeachPuff": color.RGBA{R: 255, G: 218, B: 185, A: 255},
	//MistyRose RGBA color
	"MistyRose": color.RGBA{R: 255, G: 228, B: 225, A: 255},
	//LavenderBlush RGBA color
	"LavenderBlush": color.RGBA{R: 255, G: 240, B: 245, A: 255},
	//Linen RGBA color
	"Linen": color.RGBA{R: 250, G: 240, B: 230, A: 255},
	//OldLace RGBA color
	"OldLace": color.RGBA{R: 253, G: 245, B: 230, A: 255},
	//PapayaWhip RGBA color
	"PapayaWhip": color.RGBA{R: 255, G: 239, B: 213, A: 255},
	//SeaShell RGBA color
	"SeaShell": color.RGBA{R: 255, G: 245, B: 238, A: 255},
	//MintCream RGBA color
	"MintCream": color.RGBA{R: 245, G: 255, B: 250, A: 255},
	//SlateGray RGBA color
	"SlateGray": color.RGBA{R: 112, G: 128, B: 144, A: 255},
	//LightSlateGray RGBA color
	"LightSlateGray": color.RGBA{R: 119, G: 136, B: 153, A: 255},
	//LightSteelBlue RGBA color
	"LightSteelBlue": color.RGBA{R: 176, G: 196, B: 222, A: 255},
	//Lavender RGBA color
	"Lavender": color.RGBA{R: 230, G: 230, B: 250, A: 255},
	//FloralWhite RGBA color
	"FloralWhite": color.RGBA{R: 255, G: 250, B: 240, A: 255},
	//AliceBlue RGBA color
	"AliceBlue": color.RGBA{R: 240, G: 248, B: 255, A: 255},
	//GhostWhite RGBA color
	"GhostWhite": color.RGBA{R: 248, G: 248, B: 255, A: 255},
	//Honeydew RGBA color
	"Honeydew": color.RGBA{R: 240, G: 255, B: 240, A: 255},
	//Ivory RGBA color
	"Ivory": color.RGBA{R: 255, G: 255, B: 240, A: 255},
	//Azure RGBA color
	"Azure": color.RGBA{R: 240, G: 255, B: 255, A: 255},
	//Snow RGBA color
	"Snow": color.RGBA{R: 255, G: 250, B: 250, A: 255},
	//DimGray RGBA color
	"DimGray": color.RGBA{R: 105, G: 105, B: 105, A: 255},
	//DimGrey RGBA color
	"DimGrey": color.RGBA{R: 105, G: 105, B: 105, A: 255},
	//Grey RGBA color
	"Grey": color.RGBA{R: 128, G: 128, B: 128, A: 255},
	//DarkGray RGBA color
	"DarkGray": color.RGBA{R: 169, G: 169, B: 169, A: 255},
	//DarkGrey RGBA color
	"DarkGrey": color.RGBA{R: 169, G: 169, B: 169, A: 255},
	//LightGray RGBA color
	"LightGray": color.RGBA{R: 211, G: 211, B: 211, A: 255},
	//LightGrey RGBA color
	"LightGrey": color.RGBA{R: 211, G: 211, B: 211, A: 255},
	//Gainsboro RGBA color
	"Gainsboro": color.RGBA{R: 220, G: 220, B: 220, A: 255},
	//WhiteSmoke RGBA color
	"WhiteSmoke": color.RGBA{R: 245, G: 245, B: 245, A: 255},
}
