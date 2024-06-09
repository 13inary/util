package util

func TextRed(src string) string {
	return "\033[31m" + src + "\033[0m"
}

func TextGreen(src string) string {
	return "\033[32m" + src + "\033[0m"
}

func TextYellow(src string) string {
	return "\033[33m" + src + "\033[0m"
}

func TextBlue(src string) string {
	return "\033[34m" + src + "\033[0m"
}

func TextPurple(src string) string {
	return "\033[35m" + src + "\033[0m"
}
