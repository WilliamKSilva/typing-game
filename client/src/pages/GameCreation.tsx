import axios from "axios";
import { Component, Setter, createSignal } from "solid-js";
import { Button, ButtonSize, ButtonType } from "../components/Button";
import { Input } from "../components/Input";
import "./GameCreation.css";

import { useNavigate } from "@solidjs/router";
import { FaSolidArrowLeft } from "solid-icons/fa";
import { GameState } from "../types/game-data";

type GameCreationProps = {
  setGameData: Setter<GameState>;
};

export const GameCreation: Component<GameCreationProps> = (props) => {
  const navigate = useNavigate();

  const [loading, setLoading] = createSignal(false);

  const onFormSubmit = async (event: any) => {
    event.preventDefault();
    const data = new FormData(event.target);

    const url = `${import.meta.env.VITE_APP_HTTP_SERVER_URL}/games`;
    const newGameData = {
      gameName: data.get("name"),
      playerName: data.get("player"),
    };

    try {
      setLoading(true);
      const response = await axios.post(url, newGameData);

      const responseData = response.data as GameState;

      props.setGameData(responseData);

      setLoading(false);

      navigate(`/game/${responseData.id}?player=${data.get("player")}`);
    } catch (error) {
      console.log(error);
      setLoading(false);
    }
  };

  return (
    <div class="wrapper">
      <div class="content-wrapper">
        <div class="content">
          <div class="form-header">
            <div class="arrow-left">
              <FaSolidArrowLeft size={25} color="white" onClick={() => navigate("/")} />
            </div>
            <span class="form-header-title">Create a new game</span>
          </div>
          <form class="form" onSubmit={(event) => onFormSubmit(event)}>
            <div class="form-input-wrapper">
              <Input 
                name="name" 
                placeholder="Game name" 
                disabled={false} 
              />
            </div>
            <div class="form-input-wrapper">
              <Input 
                name="player" 
                placeholder="Your nickname" 
                disabled={false} 
              />
            </div>

            <div class="form-input-button-wrapper">
              <Button
                onClick={() => {}}
                description="Salvar"
                size={ButtonSize.medium}
                type={ButtonType.submit}
                primary={true}
                loading={loading()}
              />
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};
