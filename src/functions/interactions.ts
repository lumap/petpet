import { InteractionResponseType, APIInteractionResponseCallbackData, RESTPostAPIInteractionCallbackJSONBody, APIInteraction, APIApplicationCommandInteraction, } from "discord-api-types/v10";
import { FastifyReply } from "fastify";
import axios from 'axios';
import FormData from 'form-data';

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

export function editMessage(interaction: APIInteraction, data: APIInteractionResponseCallbackData) {
    axios(`https://discord.com/api/v10/webhooks/${interaction.application_id}/${interaction.token}/messages/@original`, {
        method: 'PATCH',
        headers: {
            'Content-Type': 'application/json'
        },
        data: JSON.stringify(data)
    });
};


export async function editMessageWithAttachments(interaction: APIApplicationCommandInteraction, data: APIInteractionResponseCallbackData, buffers: Buffer[]) {
    const formData = new FormData();
    formData.append('payload_json', JSON.stringify(data));

    for (let i = 0; i < buffers.length; i++) {
        formData.append(`files[${i}]`, buffers[i], data.attachments!.find(c => c.id === i)!.filename);
    }

    axios(`https://discord.com/api/v10/webhooks/${interaction.application_id}/${interaction.token}/messages/@original`, {
        method: 'PATCH',
        headers: {
            'Content-Type': `multipart/form-data; boundary=${formData.getBoundary()}`
        },
        data: formData
    });
};