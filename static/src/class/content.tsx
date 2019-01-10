import { Place, VehicleRouting } from '../class';
import 'whatwg-fetch';

export interface ContentI {
    distances: Array<Array<number>>
    durations: Array<Array<number>>
    serviceTime: Array<number>
    demand: Array<number>
    capacity: Array<number>
    timewindows: Array<string>
    timeStart: string
}

export interface ContentToSend {

}
let limit = 5
export class Content implements ContentI{
    distances: Array<Array<number>>
    durations: Array<Array<number>>
    serviceTime: Array<number>
    posOrigin: number
    posDestination: number
    demand: Array<number>
    capacity: Array<number>
    timewindows: Array<string>
    timesStart: Array<string>
    timeStart: string
    markers: Array<google.maps.Marker>

    

    constructor(){
        this.distances = new Array()
        this.durations = new Array()
        this.serviceTime = new Array()
        this.posOrigin = 0
        this.posDestination = 0
        this.demand = new Array()
        this.capacity = new Array()
        this.capacity = [550]
        this.timewindows = new Array()
        this.timesStart = new Array()        
    }
    addDistances = () =>{
        
    }

    

    getDistancesFromCoordinates = (coordinates: Array<Array<number>>) => {
        let distancesMinutes = coordinates.map((coordinate1: Array<number>) => {
            return coordinates.map((coordinate2: Array<number>) => {
                let x = coordinate2[0]-coordinate1[0]
                let y = coordinate2[1]-coordinate1[1]
                let distance = Math.sqrt(Math.pow(x,2) + Math.pow(y,2))
                return distance
            })
        })

        this.distances = distancesMinutes.map((distances:Array<number>) => {
            return distances.map((distance: number) => {
                return distance*60
            })
        })        
        
        this.changeDistancesToDurations()
        
    }

    changeDistancesToDurations = () => {
        let speed = 1
        this.durations = this.distances.map((distancesStation: Array<number>) =>{
            return distancesStation.map((distance:number) => {
                return distance / speed
            })
        })
    }

    setTimesWindowsFromMinutes = (minutes: Array<number>) => {
       this.timewindows =  minutes.map((minute: number) => {         
            return this.minuteToTime(minute).toISOString()
        })
    }

    setTimesStartFromMinutes = (minutes:Array<number>) => {
        this.timesStart = minutes.map((minute:number) => {
            return this.minuteToTime(minute).toISOString()
        })
    }


    minuteToTime = (minute : number):Date => {
        let today = new Date()
        today = new Date(today.getFullYear(), today.getMonth(), today.getDay(),0)
        let date = new Date(today.getTime() + (minute * 60000))
        return date
    }



    updatePos = () => {

        let postEnd =  this.markers.length-1

        if(this.posDestination+limit > postEnd){
            this.posDestination = 0
            if(this.posOrigin+limit > postEnd){
                this.posOrigin = 0
            }else{
                this.posOrigin += limit
            }
        }else{
            this.posDestination += limit
        }
    }


    save = () => {
        localStorage.setItem('content', JSON.stringify(this.toString()) )
    }

    toString(){  
        return { distances: JSON.stringify(this.distances), durations: JSON.stringify(this.durations), posOrigin: this.posOrigin, posDestination: this.posDestination}
    }

    setLocalStorage = (content:any) => {
        this.distances = JSON.parse(content.distances)
        this.durations = JSON.parse(content.durations)
        this.posOrigin = content.posOrigin
        this.posDestination = content.posDestination
    }
    
    
    send = (callback: Function) => {

        this.durations = this.durations.map((durations:Array<number>, i:number) => {
            return durations.map((duration:number, j:number)=>{
                if(duration == 0 && i != j){
                    duration = 1
                }
                return duration
            })
        })
        let body = JSON.stringify({
                distances: this.durations,
                demand: this.demand,
                capacity: this.capacity,
                serviceTime: this.serviceTime,
                timewindows: this.timewindows,
                timesStart: this.timesStart,
                timeStart: this.timeStart
            })
          
            
            console.log("this.durations ", this.durations);
            

        fetch('api/ants',{
            method: 'post',
            headers:{
                'Content-Type': 'application/json'
            },
            body: body
                
        })
        .then((response:Response) =>{
            if(response.status <= 400){
                response.json().then((vehiclesRouting : Array<VehicleRouting>) => {
                    callback(vehiclesRouting)                    
                    
                    //callback()
                })
            }else{
                response.text().then((text:String) => {
                    throw text
                })
            }
        })
    }

}