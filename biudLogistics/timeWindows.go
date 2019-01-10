package biudLogistics

import(
  "github.com/gorilla/context"
  "net/http"
  "log"
)

func PostTimeWindows(w http.ResponseWriter, r *http.Request)  {
  TAG:="PostTimeWindows(w http.ResponseWriter, r *http.Request)"
  log.Println(TAG)
  var(  //  Declaracion de variables
    deposits Deposits  //  Vector de depositos
    clients Clients //  Clientes
    vehicles Vehicles //  Vehiculos
    solucions SliceRouteList //  Vector de soluciones
  )

  solucions=make(SliceRouteList,1)
  content :=context.Get(r,"body").(*Content)  // Obtiene los datos enviados del Usuario
  log.Printf("%+v\n",content)
  deposits.Get(content)  //  Obtiene los depositos enviados por el Usuario
  clients.Get(content)  //  Obtiene los clientes enviados por el Usuario
  vehicles.Get(content) //  Obtiene los vehiculos enviados por el Usuario
  log.Println("deposit,clients,vehicles",deposits,clients,vehicles)
  solucions[0].GetTimeWindows(&vehicles,&deposits,&clients,&content.TimeStart)
  solucions[0].GetTimesVisited(&deposits, &clients)
  solucions[0].LoadPaths()
  log.Println(solucions)
  responseJson(w, 202, solucions)
}
