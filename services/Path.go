package services

import (

)

func (path *Path)Retime(clientI *Station)  { //  Selecciona el nuevo tiempo de camino
  if (*path).I==(*clientI).Id {
    (*path).Time=(*clientI).Distances[(*path).J]
  }

}
