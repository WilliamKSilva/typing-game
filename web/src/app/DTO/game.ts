import { PlayerDTO } from "./player"

export type GameDTO = {
  id: string
  playerOne: PlayerDTO
  playerTwo: PlayerDTO
}
