
import "dotenv/config";

import * as joi from "joi";

interface EnvVars {
  PORT: number;
  JWT_SECRET: string;
}

const envSchema = joi
  .object<EnvVars>({
    PORT: joi.number().port().required(),
      JWT_SECRET: joi.string().required(),
  })
  .unknown(true);

const { error, value } = envSchema.validate(process.env);

if (error) {
  throw new Error(`Config validation error: ${error.message}`);
}
const envVars: EnvVars = value;
export const envs = {
  port: envVars.PORT,
  jwtSecret: envVars.JWT_SECRET
};
