package main

import (
	"context"
	"fmt"
	"sync"
)

func quarantadue(wg *sync.WaitGroup, i int){
	ctx:= context.Background()
	fmt.Println("CREATO IL BACKGROUND per quella con i = ", i)
	if(i==0){
		//wg.Done()
		return;
	}
	wg.Done()
	fmt.Println("DOPO LA RETURN QUESTO NON DOVREBBE ESSERE STAMPATO")
	<-ctx.Done()
}

func main(){
	//creo un waitgroup
	fmt.Println("INIZIO")
	var wg sync.WaitGroup //creo wg di tipo WaitGroup
	wg.Add(2) //devo aspettare 1 goroutine
	fmt.Println("ADD 2")
	go func(wg *sync.WaitGroup){
		fmt.Println("DENTRO")
		quarantadue(wg,0)//codice goroutine
		fmt.Println("FINITO 0")
	}(&wg)
	go func(wg *sync.WaitGroup){
		fmt.Println("DENTRO")
		quarantadue(wg,42)//codice goroutine
		fmt.Println("FINITO 42")
	}(&wg)
	fmt.Println("STO ASPETTANDO?")
	wg.Wait()
	fmt.Println("FINITO")
}