import { Client, GatewayIntentBits, CommandInteraction, UserContextMenuCommandInteraction, MessageContextMenuCommandInteraction, ChatInputCommandInteraction, ActivityType } from "discord.js";
let urlcheck = require('is-a-url');
const { parse } = require('twemoji-parser');
const config = require("../config.js")
const petpet = require('pet-pet-gif');
const client = new Client({ intents: [GatewayIntentBits.Guilds] });

async function getPetGif(content: string, options: any): Promise<Buffer | string> {
    let gif: any;
    try {
        gif = await petpet(content, options);
    } catch {
        gif = "ERROR";
    }
    return gif;
}


async function getSlashURL(interaction: ChatInputCommandInteraction): Promise<{ content: string, target: string }> {
    let content = "tf did u do";
    let target = "h"
    switch (interaction.options.getSubcommand()) {
        case "user": {
            const user = interaction.options.getUser("user")!;
            if (interaction.guild) {
                try {
                    const member = await interaction.guild.members.fetch(user);
                    target = member.user.tag;
                    content = member.displayAvatarURL({ extension: "png", size: 1024 });
                } catch {
                    target = user.tag;
                    content = user.displayAvatarURL({ extension: "png", size: 1024 })!;
                }
            } else {
                target = user.tag;
                content = user.displayAvatarURL({ extension: "png", size: 1024 })!;
            }
            break;
        }
        case "attachment": {
            const url = interaction.options.getAttachment("attachment")!.url;
            if (!urlcheck(url) || !isImage(url)) {
                content = "Sorry, this attachment does not seem to be valid. Please make sure it's a `jpg`, `jpeg` or `png` image.";
            } else {
                content = url;
            }
            target = "an attachment";
            break;
        }
        case "imageurl": {
            const url = interaction.options.getString("url")!;
            if (!urlcheck(url) || url.startsWith("https://tenor.com/view/") || !isImage(url)) {
                content = "Sorry, this link does not seem to be valid. Please make sure the image link ends with `.jpg`, `.jpeg` or `.png`.";
            } else {
                content = url;
            }
            target = "an image from an external URL"
            break;
        }
        case "server": {
            if (!interaction.guild) {
                content = "This command is only possible in servers";
                break;
            }
            const url = interaction.guild.iconURL({ extension: "png", size: 1024 });
            if (!url) {
                content = "This server does't have an icon."
                break;
            }
            content = url;
            target = "this server's icon"
            break;
        }
        case "emoji": {
            const emoji = interaction.options.getString("emoji")!;
            if ((emoji.match(/(<a?)?:\w+:(\d{16,20}>)?/u)) !== null) {
                content = `https://cdn.discordapp.com/emojis/${emoji.split(":")[2].slice(0, -1)}.png`;
                content = content.replaceAll(">", ""); //fixes a bug if multiple emojis are set as the argument
            } else if (parse(emoji)?.[0]?.url) {
                content = "Hello, default emojis are broken for the time being. Sorry for the interruption."
            } else {
                content = "I wasn't able to find an emoji in this. I wish discord had an \"emoji\" option for slash commands";
            }
            target = "an emoji";
            break;
        }
        default: {
            break;
        }
    }
    return { content, target };
}


function invite(interaction: ChatInputCommandInteraction) {
    try {
        return interaction.reply({
            content: "Click the button below to add me to one of your servers, or share this link to your friends to invite me: <https://lumap.fr/petpet>",
            components: [{
                type: 1,
                components: [
                    {
                        type: 2,
                        label: "Invite Me",
                        style: 5,
                        url: "https://discord.com/api/oauth2/authorize?client_id=966271016953327649&permissions=51200&scope=applications.commands%20bot"
                    }
                ]
            }],
            ephemeral: true
        });
    } catch { }
}

function support(interaction: ChatInputCommandInteraction) {
    try {
        return interaction.reply({
            content: "remind lumap to delete this kthx",
            ephemeral: true
        });
    } catch { }
}

function isImage(url: string) {
    return /\.(jpg|jpeg|png)$/.test(url.toLowerCase());
}


function sendGif(interaction: CommandInteraction, gif: Buffer, target: string) {
    interaction.editReply({
        files: [
            {
                attachment: gif,
                name: "pet.gif",
                description: `${interaction.user.tag} has pet ${target}`
            }
        ],
    }).then().catch(() => { });
}

async function handleMessageContextMenu(interaction: MessageContextMenuCommandInteraction, client: Client) {
    try { await interaction.deferReply() } catch { "why do u keep crashing here"; return; }
    let content: string = "nice try", target: string;
    if (interaction.guild) {
        try {
            if (interaction.targetMessage.webhookId) {
                target = interaction.targetMessage.author.username + " (webhook)"
                content = interaction.targetMessage.author.displayAvatarURL({ extension: "png", size: 1024 })
            } else {
                const member = await interaction.guild.members.fetch(interaction.targetMessage.author.id);
                target = member?.user?.tag;
                content = member.displayAvatarURL({ extension: "png", size: 1024 });
            }
        } catch (e: any) {
            console.log(e);
            content = "This member doesn't seem to be here. If you want to petpet them, use their user ID as the `user` argument of `/petpet user`.";
            target = "h";
        }
    } else {
        const user = await client.users.fetch(interaction.targetMessage.author.id);
        target = user.tag;
        content = user.avatarURL({ extension: "png", size: 1024 })!;
    }
    if (!content.startsWith("http")) {
        return interaction.editReply({
            content: content
        });
    }
    let options = {
        resolution: 128,
        delay: 30
    };
    const gif = await getPetGif(content, options);
    if (typeof gif === "string") {
        return interaction.editReply({ content: "I fucked up" });
    }
    sendGif(interaction, gif, target)
}


async function handleUserContextMenu(interaction: UserContextMenuCommandInteraction, client: Client) {
    try { await interaction.deferReply() } catch { "why do u keep crashing here"; return; }
    let content: string = "nice try";
    let target: string = "h";
    if (interaction.guild) {
        try {
            const member = await interaction.guild.members.fetch(interaction.targetId);
            target = member.user.tag;
            content = member.displayAvatarURL({ extension: "png", size: 1024 });
        } catch {
            content = "This member doesn't seem to be here. If you want to petpet them, use their user ID as the `user` argument of `/petpet user`.";
        }
    } else {
        const user = await client.users.fetch(interaction.targetId);
        target = user.tag;
        content = user.avatarURL({ extension: "png", size: 1024 })!;
    }
    if (!content.startsWith("http")) {
        return interaction.editReply({
            content: content
        });
    }
    let options = {
        resolution: 128,
        delay: 30
    };
    const gif = await getPetGif(content, options);
    if (typeof gif === "string") {
        return interaction.editReply({
            content: "I fucked up"
        });
    }
    sendGif(interaction, gif, target)
}

async function handleSlashCommand(interaction: ChatInputCommandInteraction) {
    const ephemeral = interaction.options.getBoolean("ephemeral") || false;
    if (interaction.commandName === "invite") {
        invite(interaction);
        return;
    } else if (interaction.commandName !== "petpet") {
        support(interaction);
        return;
    }
    try { await interaction.deferReply({ ephemeral: ephemeral }) } catch { "why do u keep crashing here"; return; }
    let { content, target } = await getSlashURL(interaction);
    if (!content?.startsWith("http")) {
        return interaction.editReply({
            content: content
        })
    }
    let gif: Buffer | string
    try {
        let options = {
            resolution: 128,
            delay: 30
        };
        if (interaction.isChatInputCommand()) {
            if (interaction.options.getInteger("delay")) {
                options.delay = interaction.options.getInteger("delay")!;
            }
            if (interaction.options.getInteger("resolution")) {
                options.resolution = interaction.options.getInteger("resolution")!;
            }
        }
        gif = await getPetGif(content, options);
        if (typeof gif === "string") { return; }
    } catch {
        interaction.editReply({
            content: "Sorry, but it looks like something went wrong. Please retry with a valid file/link",
        });
        return;
    }
    sendGif(interaction, gif!, target)
}

client.on("interactionCreate", async (interaction) => {
    if (!interaction.isCommand()) return;
    if (interaction.isUserContextMenuCommand()) {
        handleUserContextMenu(interaction, client);
    } else if (interaction.isMessageContextMenuCommand()) {
        handleMessageContextMenu(interaction, client)
    } else { //slash command
        await handleSlashCommand(interaction)
    }
})

client.login(config.token);

client.on("ready", () => {
    console.log("bot started ig")
    setActivity(client)
    setInterval(setActivity, 3600000, client)
})

function setActivity(client: Client) {
    client.user!.setPresence({ activities: [{ name: `petpets`, type: ActivityType.Watching }] });
}