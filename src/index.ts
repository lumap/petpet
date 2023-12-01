// import commands from "./commands";
import { APIInteraction, InteractionResponseType, InteractionType } from "discord-api-types/v10";
import { FastifyReply, FastifyRequest, fastify } from "fastify";
import { verify } from "discord-verify/node";
require('dotenv').config();
import crypto from "crypto";
import { handleCommands } from "./functions/handleCommands";
import { logBoot } from "./functions/logs";

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
        console.log("Invalid signature");
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
