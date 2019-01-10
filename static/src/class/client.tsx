
import {Place,Position, Deposit, User} from './'
import { isNullOrUndefined } from 'util';


export class Client extends Place {
    clientName:string
    nit:string
    comercialEstablishment:string
    address: string

    constructor(name:string,clientName:string,nit:string,comercialEstablishment:string,address:string,latLng:google.maps.LatLng,finalized:boolean = false, distances:{[key:string]:number|boolean} = {}){
        super(name,latLng,finalized,distances)

        let infoWindow = new google.maps.InfoWindow({
            content: '<div>'
                +'<div> Nombre: '+ name + '</div>'
                +"<div> Nombre del cliente: "+ clientName +"</div>"
                +"<div> NIT: "+ nit +"</div>"
                +"<div> Establecimiento Comercial:"+comercialEstablishment+"</div>"
                +"<div> Dirección:"+address+"</div>"
                +'</div>',
            position:latLng
        })

        this.clientName = clientName
        this.nit = nit
        this.comercialEstablishment = comercialEstablishment
        this.address = address

        this.infoWindow = infoWindow
    }

    toJson(){
        const {name,clientName,nit,comercialEstablishment,address,latLng,distances,finalized} = this
        return {name,clientName,nit,comercialEstablishment,address,latLng:{lat:latLng.lat(),lng:latLng.lng()},distances,finalized}
   
    }

}

export class Places {
    private limit:number;posOrigin:number;posDestination:number
    private places:(Client|Deposit)[]
    private position:number

    constructor(places:(Client|Deposit)[] = []){
        this.places = []
        this.merge(places)
        this.limit = 5
    }

    setPosition = (position:number) => {
        this.position = position
    }

    push = (place:Client|Deposit) =>{
        // console.trace("place to push: ",place.name,JSON.stringify(place.distances),place.distances)
        // console.log("place.name: ",place.name)
        if(!this.exists(place.name)){
            this.places.map((placeIn:Client|Deposit) => {
                if(isNullOrUndefined(place.distances[placeIn.name])){
                    place.distances[placeIn.name] = false                    
                }
                if(isNullOrUndefined(placeIn.distances[place.name])){
                    // console.log("entra ", placeIn.name,place.name)
                    placeIn.distances[place.name] = false
                    placeIn.finalized = false
                }                
            }) 
            // console.log("place pushed: ",place.finalized)      
            this.places.push(place)                 
        }
    }

    merge = (places:(Client|Deposit)[]) => {
        places.forEach((place:Client|Deposit) => {
            this.push(place)
        })
    }

    exists = (name:string) => {
       let placeFound =  this.getByName(name)

        return !isNullOrUndefined(placeFound)
    }

    getByName = (name:string) => {
        let placeFound:Client|Deposit
        let placesFound =  this.places.filter((place:Client|Deposit) => {
            return place.name === name
        })

        if(placesFound.length > 0){
            placeFound = placesFound[0]
        }

        return placeFound;        
    }

    toJson(){
       return this.places.map((place:Client|Deposit) => {
            return place.toJson()
        })
    }

    clientsToJson = () => {
        return this.getClients().map((client:Client) => {
            return client.toJson()
        })
    }

    depositsToJson = () => {
        return this.getDeposits().map((deposit:Deposit) => {
            return deposit.toJson()
        })
    }

    get = () => {
        return this.places
    }

    getClients = () => {
        return this.places.filter((place:Client|Deposit) => {
            return place instanceof Client
        })
    }

    getDeposits = () => {
        return this.places.filter((place:Client|Deposit) => {
            return place instanceof Deposit
        })
    }

    getOriginsAndDestinations = () => {
        const limit = 5; 
        let i = 0;
        let origins:(Client|Deposit)[] = [];
        let destinations:(Client|Deposit)[] = [];

        origins = this.places.filter((origin:Client|Deposit) => {
            if(origin.finalized == false && i < 5){
                i++
                return true;
            } else{
                return false;
            }           
        })

        i = 0
        origins.forEach((place:Client|Deposit) => {
            let distances = place.distances
            // console.log("place: ",place)
            let finalized:boolean = true
            for(const name in distances){
                
                if(distances.hasOwnProperty(name) && distances[name] === false && i < 5){

                    console.log("distances[name]: ",place.name,name,distances[name])
                    let destination = this.getByName(name)
                    destinations.push(destination)
                    // finalized = false
                    i++
                }else if(distances[name] === false ){
                    finalized = false
                }
            }
            place.finalized = finalized
        })

        return {origins:origins.slice(),destinations:destinations.slice()}
    }

    getDistances = (callback: Function, user:User) => {
        // console.log("getDistances")
        let service = new google.maps.DistanceMatrixService()
        const {origins,destinations} = this.getOriginsAndDestinations()
        const originsCoordinates = origins.map((place:Client|Deposit) => { return place.latLng})
        const destinationsCoordinates = destinations.map((place:Client|Deposit) => { return place.latLng})
        // console.log("origins: ",origins)
        // console.log("destinations: ",destinations)
        if (origins.length > 0 && destinations.length > 0) {            
            const distances = localStorage.getItem("distances")

            if(!isNullOrUndefined(distances)){
                this.readDistances(JSON.parse(distances),callback,user)
            }else{
                service.getDistanceMatrix({
                    origins: originsCoordinates,
                    destinations: destinationsCoordinates,
                    travelMode: google.maps.TravelMode.DRIVING,
                    }, 
                    (response:google.maps.DistanceMatrixResponse, status: google.maps.DistanceMatrixStatus) => {
                        this.readDistanceMatrix(response, status, origins, destinations, callback,user)
                    }
                )
            } 
            
        }else{ 
            callback(user)
        }
    }

    readDistances = (distances:{[key:string]:number[]},callback:Function,user:User) => {
        console.log("readDistances")
        for(var name in distances){
            if(distances.hasOwnProperty(name)){
                let place = this.getByName(name)
                let distanceVector = distances[name]
                if(!isNullOrUndefined(place) && !isNullOrUndefined(distanceVector)){
                    let i = 0
                    place.finalized = true
                    // console.log("name: ",name," - ",place.distances)
                    for(var key in place.distances){
                        if(place.distances.hasOwnProperty(key)){
                            // console.log("key: ",key, " - i: ",i)
                            const distance = distanceVector[i]
                            if(!isNullOrUndefined(distance)){
                                place.distances[this.places[i].name] = distanceVector[i]                                
                            }
                            i++
                        }
                    }
                }                
            }
        }
        localStorage.removeItem("distances")
        callback(user)

    }

    readDistanceMatrix = (response:google.maps.DistanceMatrixResponse, status: google.maps.DistanceMatrixStatus,origins:(Client|Deposit)[], destinations:(Client|Deposit)[], callback: Function,user:User) => {
        console.log("staus ",status)
        if (status === google.maps.DistanceMatrixStatus.OK){
            let rows = response.rows
            rows.map((row: google.maps.DistanceMatrixResponseRow, positionOfOrigin: number) => {
                let elements = row.elements
                const name = origins[positionOfOrigin].name
                let origin = this.getByName(name)
                
                // if (this.distances[posI] === undefined){ this.distances.push([])}
                // if (this.durations[posI] === undefined){ this.durations.push([])}
                elements.map((responseElement: google.maps.DistanceMatrixResponseElement, positionOfDestination:number) => {
                   if(responseElement.status == google.maps.DistanceMatrixElementStatus.OK){
                       
                        let destination = destinations[positionOfDestination]
                        origin.distances[destination.name] = responseElement.distance.value

                        if(!isNullOrUndefined(user)){
                            user.places.getByName(name).distances[destination.name] = responseElement.distance.value
                        }
                        // this.distances[posI].push(responseElement.distance.value)
                        // this.durations[posI].push(responseElement.duration.value)                    

                   }

                })
            })


            // this.save()
            // callback(user) 
            setTimeout(()=>{this.getDistances(callback,user)}, 3000)

        }else{
            console.log(status)
            setTimeout(()=>{this.getDistances(callback,user)}, 10000)
        }
    }

    save = () => {
        localStorage.setItem('clients', JSON.stringify(this.clientsToJson()) )
        localStorage.setItem("deposits", JSON.stringify(this.depositsToJson()))
    }

    getDistancesToSend = () => {
        const distances = this.places.map((clientI:Client|Deposit) => {
            // return {name:clientI.name, distance: this.places.map((clientJ:|Deposit) => {
            return  this.places.map((clientJ:|Deposit) => {
                
                // return {name:clientJ.name,distance:clientI.distances[clientJ.name]}
                return clientI.distances[clientJ.name]
                

            })
        })
        console.log("distances: ",distances)
        return distances
        

    }

    delete = (name:string) => {
        this.places = this.places.filter((place:Client|Deposit) => {
            delete place.distances[name];            
            return place.name !== name
        })
    }

}