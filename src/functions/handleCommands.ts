import { APIApplicationCommandInteraction, APIMessageApplicationCommandInteractionData, ApplicationCommandType } from "discord-api-types/v10";
import { FastifyReply } from "fastify";
import commands from "../commands/";
import { logError } from "./logs";

export async function handleCommands(interaction: APIApplicationCommandInteraction, res: FastifyReply) {
    try {
        switch (interaction.data.type) {
            case ApplicationCommandType.ChatInput: {
                switch (interaction.data.name) {
                    case "support": {
                        return commands.support(interaction.data, res);
                    }
                    case "invite": {
                        return commands.invite(interaction.data, res);
                    }
                    case "petpet": {
                        return commands.petpet(interaction.data, res, interaction);
                    }
                }
            }
            case ApplicationCommandType.Message: {
                return commands.petpetMsgCtx(interaction.data as APIMessageApplicationCommandInteractionData, res, interaction);
            }
            case ApplicationCommandType.User: {
                return commands.petpetUserCtx(interaction.data, res, interaction);
            }
        };
    } catch (e) {
        logError(e);
    }
}