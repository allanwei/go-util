package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/allanwei/gcwebapis-db"
	//"github.com/allanwei/gcwebapis-util"
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args
	t := args[1]
	n, err := strconv.Atoi(args[2])
	if err != nil {
		log.Fatal(err)
	}
	f, err := strconv.Atoi(args[3])
	if err != nil {
		log.Fatal(err)
	}
	if err = insertTo(t, n, f); err != nil {
		log.Fatal(err)
	}

}
func randint() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2)

}
func randfloat() float64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()

}
func insertTo(t string, n, f int) error {
	con, err := db.CreateDBCon(nil)
	if err != nil {
		return err
	}
	defer con.Close()
	sqlexpress := fmt.Sprintf(`INSERT INTO [dbo].[Data_Table]([Time_Stamp],[GMX_1],[RPM_1],[RTD_2],[RTD_3],
		[CTX_1],[CTX_2],[CTX_3],[CTX_4],
		[POT_1],[POT_2],[POT_3],
		[PTX_1],[PTX_2],[PTX_3],[PTX_4],[PTX_5],[PTX_6],[PTX_7],[PTX_8],[PTX_9],[PTX_10],[PTX_11],[PTX_12],
		[EPB_1],[EPB_2],[EXT_1],[Torque],[Water],[Soap],[Air],[ROA],[Thrust],
		[PIPE_No],[FER],[FIR],
		[LSW_5],[LSW_6],[LSW_7A],[LSW_7B],[LSW_8],[LSW_9A],[LSW_9B],[LSW_10],[LSW_11],[LSW_12A],[LSW_12B],[PSL_1_L],
		[PSL_1_H],[PSL_2_L],[PSL_2_H],[PSL_3],[PSL_4],[PSL_5],[PSL_6],[PSL_7],[PSL_8],[PIPE_PUSH])
	VALUES('%s',%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%d,%f,
	%f,%f,%f,%f,%d,%f,%f,%f,%f,%f,%f,%d,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%d)`, t, randfloat(), randfloat(), randfloat(), randfloat(), randfloat(),
		randfloat(), randfloat(), randfloat(), randfloat(), randfloat(), randfloat(),
		randfloat(), randfloat(), randfloat(), randfloat(), randfloat(), randfloat(), randfloat(),
		randfloat(), randfloat(), randfloat(), randfloat(), randfloat(), randfloat(), randfloat(), randfloat(),
		randfloat(), randfloat(), randfloat(), randfloat(), randfloat(), randfloat(), n, randfloat(), randfloat(),
		randfloat(), randfloat(), randfloat(), randint(), randfloat(), randfloat(), randfloat(), randfloat(), randfloat(),
		randfloat(), randint(), randfloat(), randfloat(), randfloat(), randfloat(), randfloat(), randfloat(), randfloat(),
		randfloat(), randfloat(), randfloat(), f)

	_, err = con.ExecuteQuery(sqlexpress)
	if err != nil {
		return err
	}
	//runtime.Goexit()
	return nil
}
