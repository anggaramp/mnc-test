package main

import (
	"fmt"
	"time"
)

func main() {
	const totalLeave = 14
	const layout = "2006-01-02"

	var totalMassLeave, leaveDuration int
	var joinDate, leaveDate string
	fmt.Print("Jumlah Cuti Bersama    : ")
	fmt.Scan(&totalMassLeave)
	fmt.Print("Tanggal join karyawan  : ")
	fmt.Scan(&joinDate)
	fmt.Print("Tanggal rencana cuti   : ")
	fmt.Scan(&leaveDate)
	fmt.Print("Durasi cuti (hari)     : ")
	fmt.Scan(&leaveDuration)

	// Parse the start and end dates
	startJoin, err := time.Parse(layout, joinDate)
	if err != nil {
		fmt.Println("Error parsing start date:", err)
		return
	}
	startLeave, err := time.Parse(layout, leaveDate)
	if err != nil {
		fmt.Println("Error parsing start date:", err)
		return
	}

	endDate := time.Date(startJoin.Year(), 12, 31, 0, 0, 0, 0, startJoin.Location())

	totalDay := endDate.Sub(startJoin.Add(time.Hour*24*180)).Hours() / 24

	personalLeave := totalLeave - totalMassLeave

	personalLeave = int(totalDay / 365 * float64(personalLeave))

	if leaveDuration > 3 {
		fmt.Println(false)
		return
	}

	if startLeave.Before(startJoin.Add(time.Hour * 24 * 180)) {
		fmt.Println(false)
		fmt.Println("Alasan : Karena belum 180 hari sejak tanggal join karyawan")
		return
	}
	if leaveDuration > personalLeave {
		fmt.Println(false)
		fmt.Println(fmt.Sprintf("Alasan: Karena hanya boleh mengambil %d hari cuti", personalLeave))
		return
	}
	fmt.Println(true)
	return
}
