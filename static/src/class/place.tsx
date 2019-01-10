export class Place {
    name: string
    description?: string
    latLng: google.maps.LatLng
    infoWindow:google.maps.InfoWindow
    distances:{[key:string]:number|boolean}
    finalized:boolean

    constructor(name:string,latLng: google.maps.LatLng,finalized:boolean = false, distances:{[key:string]:number|boolean} = {}){

        let infoWindow = new google.maps.InfoWindow({
            content: '<div>'
                +'<div> Name: '+ name + '</div>'
                +'</div>',
            position:latLng
        })

        this.name = name

        this.latLng = latLng
        this.infoWindow = infoWindow     
        this.distances = distances
        this.distances[name] = 0;
        this.finalized = finalized;
    }

    toJson(){
        const {name,latLng,distances,finalized} = this
        return {name,latLng:{lat:latLng.lat(),lng:latLng.lng()},distances,finalized}
    }

    
}