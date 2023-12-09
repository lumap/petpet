import axios from 'axios';
import { APIApplicationCommandInteraction, APIInteraction, APIInteractionResponseCallbackData, InteractionResponseType, RESTPostAPIInteractionCallbackJSONBody, } from "discord-api-types/v10";
import { FastifyReply } from "fastify";
import FormData from 'form-data';
import { logError } from "./logs";

export function sendMessage(res: FastifyReply, data: APIInteractionResponseCallbackData) {
    return res.code(200).send({
        type: InteractionResponseType.ChannelMessageWithSource,
        data: data
    } as RESTPostAPIInteractionCallbackJSONBody);
}

export function deferInteration(res: FastifyReply) {
    res.code(200).send({
        type: InteractionResponseType.DeferredChannelMessageWithSource
    } as RESTPostAPIInteractionCallbackJSONBody);
}

export async function editMessage(interaction: APIInteraction, data: APIInteractionResponseCallbackData) {
    try {
        await axios(`https://discord.com/api/v10/webhooks/${interaction.application_id}/${interaction.token}/messages/@original`, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json'
            },
            data: JSON.stringify(data)
        });
    } catch (e) {
        logError(e);
    }
};


export async function editMessageWithAttachments(interaction: APIApplicationCommandInteraction, data: APIInteractionResponseCallbackData, buffers: Buffer[]) {
    const formData = new FormData();
    formData.append('payload_json', JSON.stringify(data));

    for (let i = 0; i < buffers.length; i++) {
        formData.append(`files[${i}]`, buffers[i], data.attachments!.find(c => c.id === i)!.filename);
    }

    try {
        await axios(`https://discord.com/api/v10/webhooks/${interaction.application_id}/${interaction.token}/messages/@original`, {
            method: 'PATCH',
            headers: {
                'Content-Type': `multipart/form-data; boundary=${formData.getBoundary()}`
            },
            data: formData
        });
    } catch (e) {
        logError(e);
    }
};