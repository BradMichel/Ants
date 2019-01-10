import {Client,Deposit,Places, VehicleRouting, Place} from "./"
import firebase, { FirebaseFirestore } from "../config/firebase"
import { firestore } from "firebase/app";

export class User {
    name:string
    ID:string
    places:Places

    constructor(name:string,ID:string,clients:Array<Client>,deposits:Array<Deposit>){
        this.name=name
        this.ID=ID
        this.places= new Places(deposits)
        this.places.merge(clients)

    }

    create = (callback:(user:User)=>void) => {

        const create = confirm("Esta seguro de crear un nuevo usuario")
        if(create){
            FirebaseFirestore().collection("users").add({
                name: this.name,
                clients: [],
                deposits: []
            })
            .then((docRef:firebase.firestore.DocumentReference) => {
                this.ID = docRef.id
                callback(this)
            })
            .catch((error:firebase.firestore.FirestoreError) => {
                throw error;
            })
        }else{
            alert("Usuario no creado")
        }  

    }

    saveClients = () => {
        // console.trace("saveClients")
        
        let clients = this.places.clientsToJson()
        // console.log("ID",this.ID,"saveClients: ",JSON.stringify(clients)) 

        FirebaseFirestore().collection("users").doc(this.ID).update({
            clients: clients
        })
        .then(() => {
            // console.log("-----------------value")
        })
        .catch((error:firebase.firestore.FirestoreError) => {
            throw error
        })
    }

    saveDeposits = () => {
        // console.trace("saveDeposits: ", this.places.depositsToJson())
        FirebaseFirestore().collection("users").doc(this.ID).update({
            deposits: this.places.depositsToJson()
        })
        .catch((error:firebase.firestore.FirestoreError) => {
            throw error
        })
    }

    get = () => {
        FirebaseFirestore().collection("users").doc(this.ID).get()
        .then((documentSnapshot:firebase.firestore.DocumentSnapshot) => {
            let data = documentSnapshot.data()

            this.name = data.name
            this.places = new Places(data.deposits)
            this.places.merge(data.clients)            
        })
    }

    send = (places:Places, callback:Function) => {

        const body = JSON.stringify({
            distances: places.getDistancesToSend(),
            demand: places.get().map((place:Client|Deposit) => { const demand = (place instanceof Deposit) ? 0 : 1 ; return demand;}),
            capacity: [places.get().length*2]
        })    

        console.log("body: ",body)

        fetch(`api/ants`,{
            method:`post`,
            headers:{
                'Content-Type': 'application/json'
            },
            body:body  
        })
        .then((response:Response) => {
            if(response.status <= 400){
                response.json().then((vehiclesRouting:VehicleRouting[]) => {
                    if(vehiclesRouting.length > 0){
                        console.log("places: ",places.get())
                        console.log("route: ",vehiclesRouting[0].path)
                        callback(vehiclesRouting[0])
                    }
                    
                })
            }else{
                response.text().then((text:String) => {
                    throw text
                })
            }
        })
    }

    getDistances = () => {
        this.places.getDistances((user:User) => {user.saveClients();user.saveDeposits()},this)
    }
    
}

export class Users{
    private users:Array<User>
    
    constructor(){
        this.users = new Array()
    }

    getUser = (ID:string) => {
       const usersFound = this.users.filter((user:User)=>{return user.ID === ID})
       const user = (usersFound.length > 0) ? usersFound[0]:undefined
       return user
    }

    load = () => {
        FirebaseFirestore().collection("users").get()
        .then((querySnapshot:firebase.firestore.QuerySnapshot) => {
            querySnapshot.forEach((doc:firebase.firestore.DocumentSnapshot)=>{
                const data = doc.data()
                // console.log(data.name,doc.id,data.clients)
                const user = new User(data.name,doc.id,this.getClientsFromData(data.clients),this.getDepostisFromData(data.deposits))
                this.add(user)


            })
        })



    }

    getClientsFromData = (clientsData:any) => {
        return clientsData.map((client:any)=>{ return new Client(client.name,client.clientName,client.nit,client.comercialEstablishment,client.address,new google.maps.LatLng(client.latLng.lat,client.latLng.lng),client.finalized,client.distances)})
    }

    getDepostisFromData = (depositsData:any) => {
        return depositsData.map((deposit:any) => { return new Deposit(deposit.name, new google.maps.LatLng(deposit.latLng.lat,deposit.latLng.lng),deposit.finalized,deposit.distances)})
    }

    add = (user:User)=>{
        this.users.push(user)
    }

    search = (value:string) => {
        const inputValue = value.trim().toLowerCase();
        const inputLength = inputValue.length;
        let usersFound:Array<User> =  inputLength === 0 ? [] :this.users.filter(user =>
            user.name.toLowerCase().slice(0, inputLength) === inputValue
          );

        return usersFound
    }

}