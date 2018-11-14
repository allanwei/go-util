package main

import (
	"sync"
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"github.com/allanwei/gcwebapis-util"
	//"github.com/allanwei/gocron"
	//"fmt"
	"strconv"
	"time"
	"runtime"
	"runtime/debug"
	
	"github.com/kardianos/service"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
	"path/filepath"
)
var slogger service.Logger
type program struct{}
//Start ...
func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
//Stop ...
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}
func (p *program) run() {
	run()
}
func loadConfig() (delay string,executable string,log string,currentNum int) {
	
	delay ="60s"
	temp := "go-insert.exe"
	log = "c:/temp/eum.log"
	ex, err := os.Executable()
	if err != nil {
		return
	}
	exPath := filepath.Dir(ex)
	expath := filepath.ToSlash(exPath)
	executable = fmt.Sprintf("%s/%s",expath,temp)

	err = config.Load(file.NewSource(
		file.WithPath(exPath + "/conf.json")))
	if err != nil {
		return
	}
	_ = config.Get("delay").Scan(&delay)
	_ = config.Get("log").Scan(&log)
	err = config.Get("executable").Scan(&temp)
	if err != nil{
		return
	}
	_ = config.Get("CurrentNum").Scan(&currentNum)
	executable = fmt.Sprintf("%s/%s",expath,temp)
	return 

}
type csvline struct {
	System     string
	SystenDesc string
	TagDesc    string
	TagName    string
	UOM        string
	IsUsed     string
	MachineID  string
}
func work(s  time.Time,r chan int,c chan int,execfile string,logfile string){

	//log.Println("Job Start at",time.Now())
	starttime := s
	rnumber:= <-r
	count:= <-c
	base :=float64(60)
	t :=time.Now()
	d := t.Sub(starttime)
	compare := base*float64(count)
	pushflag :=1
	n:= util.GetSQLTimeString(t)
	addnew:=false
	var wg sync.WaitGroup
	//log.Print(d.Minutes())
	if d.Minutes()>=compare{
		pushflag=0
		addnew=true
	}
	wg.Add(1)
	go func(){
		defer recovery(logfile)
		defer wg.Done()
		writeTolog(logfile,fmt.Sprintf("current pipepush= %d; time=%s; pushflag=%d; count=%d",rnumber,n,pushflag,count))
		ir:= strconv.Itoa(rnumber)
		ipf:=strconv.Itoa(pushflag)
		cmd := exec.Command(execfile,n,ir,ipf)
		if err := cmd.Start();err !=nil{
			writeTolog(logfile,err.Error())
		}
		done :=make(chan error, 1)
		go func(){
			done <-cmd.Wait()
		}()
		select{
		case <-time.After(3* time.Second):
			if err := cmd.Process.Kill(); err != nil {
				writeTolog(logfile,err.Error())
			}
			writeTolog(logfile,"process killed as timeout reached")
			log.Println("process killed as timeout reached")
		case err := <-done:
			if err != nil {
				writeTolog(logfile,err.Error())
			}
			writeTolog(logfile,"process finished successfully")
		}


		runtime.Goexit()

	}()
	wg.Wait()
	if(addnew){
		rnumber++
		count++
	}
	pushflag=1
	
	r <- rnumber
	c <- count

	log.Println("Job End at",time.Now())
	return


}
func writeTolog(logfile,str string) {
	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {

		defer f.Close()
		wr := fmt.Sprintf("%s\n", str)
		f.Write([]byte(wr))
	}

}
func startwork(d time.Duration,i int,execfile,logfile string){
	
	//now := time.Now()
	s :=time.Now()
	c :=make(chan int,1)
	r :=make(chan int,1)
	c <- 1
	r <- i
	
	for {
		timer := time.NewTimer(d)	
		

		for{
			select{
			case  <- timer.C:
				var wg sync.WaitGroup
				wg.Add(1)
				go func(){
					defer wg.Done()
					work(s,r,c,execfile,logfile)
					runtime.Goexit()
				}()
				wg.Wait()				
				
				var mem runtime.MemStats
				runtime.ReadMemStats(&mem)

				writeTolog(logfile, fmt.Sprintf("current: %fMB. Number of goroutines: %d", float32(mem.Alloc)/1024.0/1024.0, runtime.NumGoroutine()))
			}
			break
		}
	}
	
	

}

func recovery(logfile string) {  
    if r := recover(); r != nil {
		msg :=fmt.Sprintln("recovered:", r)
		writeTolog(logfile,msg)
		stack := debug.Stack()
		writeTolog(logfile,string(stack))
    }
}


//ScheduleTest ...
func ScheduleTest(d time.Duration,i int,execfile,logfile string) {
	defer recovery(logfile)
	startwork(d,i,execfile,logfile)
}
func readcsv(){
	csvFile, _ := os.Open("columns_mapping.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var csvlines []csvline
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		csvlines = append(csvlines, csvline{
			System:     line[0],
			SystenDesc: line[1],
			TagDesc:    line[2],
			TagName:    line[3],
			UOM:        line[4],
			IsUsed:     line[5],
			MachineID:  line[6],
		})
	}
	for _, l := range csvlines {
		sqlexpress := fmt.Sprintf(`INSERT INTO [dbo].[System_Tag]([SystemDesc],[SystemName],[TagDesc],[TagName],[TagUom],[IsUsed],[MachineId])
		values('%s','%s','%s','%s','%s',%s,%s)`, l.System, l.SystenDesc, l.TagDesc, l.TagName, l.UOM, l.IsUsed, l.MachineID)
		f, _ := os.OpenFile("system_tag_insert.sql", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		wr := fmt.Sprintf("%s\n", sqlexpress)
		_, _ = f.Write([]byte(wr))
		_ = f.Close()

	}
}
func run(){
	delay,execfile,logfile,num:=loadConfig()
	duration, err := time.ParseDuration(delay)
	if err != nil {
		duration = 10 * time.Second
	}
	ScheduleTest(duration,num,execfile,logfile)
}
func main() {
	_,_,logfile,_:=loadConfig()
	svcConfig :=&service.Config{
		Name : "GroundCastemulator",
		DisplayName : "GroundCast emulator",
		Description : "GroundCast emulator.",
	}
	prg := &program{}
	s,err := service.New(prg,svcConfig)
	if err !=nil{
		writeTolog(logfile,err.Error())
	}
	err = s.Run()
	if err != nil{
		writeTolog(logfile,err.Error())
	}
	
}
