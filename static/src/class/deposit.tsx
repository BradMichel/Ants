import {Place, Position} from './'
export class Deposit extends Place {

    constructor(name: string, latLng:google.maps.LatLng, finalized:boolean=false, distances:{[key:string]:number|boolean} = {}){
        super(name,latLng,finalized,distances)
    }



}
