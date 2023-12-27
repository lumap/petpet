// import commands from "./commands";
import crypto from "crypto";
import { APIInteraction, InteractionResponseType, InteractionType } from "discord-api-types/v10";
import { verify } from "discord-verify/node";
import { FastifyReply, FastifyRequest, fastify } from "fastify";
import { handleCommands } from "./functions/handleCommands";
import { logBoot } from "./functions/logs";
require('dotenv').config();

const app = fastify({ logger: false, trustProxy: 1 });

app.get('/', async (req, res) => {
    res.send("Hi there");
});


app.post('/', async (req: FastifyRequest<{
    Body: APIInteraction;
    Headers: {
        "x-signature-ed25519": string;
        "x-signature-timestamp": string;
    };
}>,
    res: FastifyReply
) => {
    const signature = req.headers["x-signature-ed25519"];
    const timestamp = req.headers["x-signature-timestamp"];
    const rawBody = JSON.stringify(req.body);

    const isValid = await verify(
        rawBody,
        signature,
        timestamp,
        process.env.PUBLICKEY!,
        crypto.webcrypto.subtle
    );

    if (!isValid) {
        return res.code(401).send("Invalid signature");
    }

    const interaction = req.body;

    switch (interaction.type) {
        case InteractionType.Ping:
            return res.send({ type: InteractionResponseType.Pong });
        case InteractionType.ApplicationCommand: {
            return handleCommands(interaction, res);
        }
        default:
            return res.code(401).send("tf did u do");
    }
});

app.listen({ port: Number(process.env.PORT), host: '0.0.0.0' }, (err, address) => {
    if (err) throw err;
    console.log("yay");
    logBoot();
});
