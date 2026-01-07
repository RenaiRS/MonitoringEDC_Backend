package utils

func GetNamaNasabah(nama, sumberData string) string {
	if sumberData == "vendor" {
		return "Belum Terdata"
	}
	return nama
}
