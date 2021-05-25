package main

import (
	"container/heap"
	"database/sql"
	"fmt"
	"math"

	_ "github.com/go-sql-driver/mysql"
)
type element struct{
	id int
	lat, lon float64
}

type Item struct {
	id int
	priority float64    
	index int 
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  
	item.index = -1 
	*pq = old[0 : n-1]
	return item
}



func distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64) float64 {
	const PI float64 = 3.141592653589793
	
	radlat1 := float64(PI * lat1 / 180)
	radlat2 := float64(PI * lat2 / 180)
	
	theta := float64(lng1 - lng2)
	radtheta := float64(PI * theta / 180)
	
	dist := 6471.01*(math.Sin(radlat1) * math.Sin(radlat2) + math.Cos(radlat1) * math.Cos(radlat2) * math.Cos(radtheta))
	
	return dist
}

func main() {
	

	db, err := sql.Open("mysql", "root:random1!@tcp(127.0.0.1:3306)/project")

	if err != nil {
        panic(err.Error())
    }

	elements, err := db.Query("SELECT * FROM vehicles")

    if err != nil {
        panic(err.Error())
    }
	
	var lat1, lon1 float64
	fmt.Println("Enter latitude and longtitude of the current user location ")

    fmt.Scanln(&lat1,&lon1)
	pq := make(PriorityQueue, 5)
	i:=0
	for elements.Next() {


		var cab element
		err := elements.Scan(&cab.id, &cab.lat,&cab.lon)
		if err != nil {
		  panic(err.Error())
		}
		

		if i<5{
			pq[i]=&Item{id:cab.id,
				priority: distance(lat1,lon1,cab.lat,cab.lon),
				index:   i,}
			i++
			if i==5{
				heap.Init(&pq)
			 	heap.Fix(&pq,1)
				heap.Fix(&pq,0)
			}
		}else{
			
			item:=&Item{id:cab.id,
				priority: distance(lat1,lon1,cab.lat,cab.lon),
				}
			if pq[4].priority>item.priority{
				heap.Pop(&pq)
				heap.Push(&pq, item)
				pq[4].id=item.id
				pq[4].priority=item.priority
				heap.Fix(&pq,1)
				heap.Fix(&pq,0)
			}	
			
		}
	}
	i=0
	for i < 5 {
		fmt.Printf("%d\n", pq[i].id)
		i++
	}
			
    
    defer db.Close()
	
	
}