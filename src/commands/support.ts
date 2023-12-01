import { FastifyReply } from "fastify";
import { sendMessage } from "../functions/interactions";
import { APIChatInputApplicationCommandInteractionData } from "discord-api-types/v10";
import { logCommand } from "../functions/logs";

export function support(interaction: APIChatInputApplicationCommandInteractionData, res: FastifyReply) {
    logCommand("support");
    return sendMessage(res, {
        content: "Join my support server through this invite: https://discord.gg/S5yryjRuse",
        flags: 64
    });
}