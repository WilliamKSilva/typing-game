import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Component } from '@angular/core';
import { GameDTO } from '../../DTO/game';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './home.component.html',
  styleUrl: './home.component.css'
})
export class HomeComponent {
  constructor(private http: HttpClient) {}
  private baseURL: string = "http://127.0.0.1:8080"

  game: GameDTO | undefined

  joinGamePayloadData = {
    gameId: '',
    playerName: ''
  }

  showGameMenu = true
  showCreateGameMenu = false
  showJoinGameMenu = false

  onInputNameChange(event: any) {
    this.joinGamePayloadData.playerName = event.target.value
  }

  onInputIDChange(event: any) {
    this.joinGamePayloadData.gameId = event.target.value
  }

  onCreateGame() {
    this.showGameMenu = false
    this.showCreateGameMenu = true
  }

  onJoinGame() {
    this.showGameMenu = false
    this.showJoinGameMenu = true
  }

  async onCreateGameConfim() {
    this.showCreateGameMenu = false
    const createGameURL: string = this.baseURL + "/create-game"
    const response = await this.http.post(createGameURL, null).toPromise() as GameDTO

    this.joinGamePayloadData.gameId = response.id

    await this.onJoinGameConfim()
  }

  async onJoinGameConfim() {
    this.showJoinGameMenu = false
    const joinGameURL: string = this.baseURL + "/join-game"
    console.log(this.joinGamePayloadData)
    const response = await this.http.post(joinGameURL, this.joinGamePayloadData).toPromise() as GameDTO

    this.game = response
  }
}
