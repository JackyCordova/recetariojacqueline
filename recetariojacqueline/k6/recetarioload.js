import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';

const client = new grpc.Client();

// Cargar proto desde src/api
client.load(['/Users/jacky/Documents/CompDistribuido/recetariojacqueline/src/api'], 'recetario.proto');


export const options = {
  stages: [
    { duration: '15s', target: 50 },
    { duration: '15s', target: 100 },
    { duration: '15s', target: 200 },
    { duration: '15s', target: 250 },
    { duration: '20s', target: 300 },
  ],
};


export default () => {
  client.connect(__ENV.RECETARIO_ADDR || '127.0.0.1:50051', {
    plaintext: true,
  });

  const data = { recipe_id: 'r1' }; // usa una receta válida

  const res = client.invoke(
    'recetario.RecipeService/GetRecipeDetails',
    data,
  );

  check(res, {
    'status OK': (r) => r && r.status === grpc.StatusOK,
  });

  client.close();
  sleep(1);
};