import app from './app';
import { envs } from './config/env.validation';

const PORT = envs.port || 4000;

app.listen(PORT, '0.0.0.0', () => {
});