import { Client, GatewayIntentBits, CommandInteraction, UserContextMenuCommandInteraction, MessageContextMenuCommandInteraction, ChatInputCommandInteraction, ActivityType } from "discord.js";
let urlcheck = require('is-a-url');
const { parse } = require('twemoji-parser');
const config = require("../config.js")
const petpet = require('pet-pet-gif');
const client = new Client({ intents: [GatewayIntentBits.Guilds], ws: { properties: { browser: "Discord iOS" } } });
let petCounter = 0;
let rateLimits: { time: number, id: string }[] = [];
const axios = require('axios').default;

(async () => {


    setInterval(function () {
        rateLimits = rateLimits.filter(c => c.time + 60000 > Date.now())
    }, 60000)


    function checkURL(interaction: ChatInputCommandInteraction): string {
        let content;
        const url = interaction.options.getString("url")!;
        if (!urlcheck(url) || url.startsWith("https://tenor.com/view/") || !isImage(url)) {
            content = "Sorry, this link does not seem to be valid. Please make sure the image link ends with `.jpg`, `.jpeg` or `.png`.";
        } else {
            content = url;
        }
        return content;
    }


    async function createGif(interaction: ChatInputCommandInteraction | MessageContextMenuCommandInteraction | UserContextMenuCommandInteraction, content: string): Promise<Buffer | boolean> {
        let gif: Buffer | string;
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
            gif = await getPetGif(content);
            if (typeof gif === "string") { return false; }
            rateLimits.push({ id: interaction.user.id, time: Date.now() });
        } catch {
            interaction.editReply({
                content: "Sorry, but it looks like something went wrong. Please retry with a valid file/link",
            });
            return false;
        }
        return gif;
    }

    function getAttachment(interaction: ChatInputCommandInteraction | MessageContextMenuCommandInteraction | UserContextMenuCommandInteraction): string {
        let content: string;
        const url = interaction.options.getAttachment("attachment")!.url;
        if (!urlcheck(url) || !isImage(url)) {
            content = "Sorry, this attachment does not seem to be valid. Please make sure it's a `jpg`, `jpeg` or `png` image.";
        } else {
            content = url;
        }
        return content;
    }

    function getEmoji(interaction: ChatInputCommandInteraction) {
        let content;
        const emoji = interaction.options.getString("emoji")!;
        if ((emoji.match(/(<a?)?:\w+:(\d{16,20}>)?/u)) !== null) {
            content = `https://cdn.discordapp.com/emojis/${emoji.split(":")[2].slice(0, -1)}.png`;
        } else if (parse(emoji)?.[0]?.url) {
            content = "Hello, default emojis are broken for the time being. Sorry for the interruption."
        } else {
            content = "I wasn't able to find an emoji in this. I wish discord had an \"emoji\" option for slash commands";
        }
        return content;
    }

    async function getPetGif(content: string): Promise<Buffer | string> {
        let gif: any;
        try {
            gif = await petpet(content);
        } catch {
            gif = "ERROR";
        }
        return gif;
    }


    async function getSlashURL(interaction: ChatInputCommandInteraction): Promise<string> {
        let content: string = "tf did u do";
        switch (interaction.options.getSubcommand()) {
            case "user": {
                content = await getUserPFP(interaction);
                break;
            }
            case "attachment": {
                content = getAttachment(interaction);
                break;
            }
            case "imageurl": {
                content = checkURL(interaction);
                break;
            }
            case "emoji": {
                content = getEmoji(interaction);
                break;
            }
            default: {
                break;
            }
        }
        return content;
    }


    async function getUserPFP(interaction: ChatInputCommandInteraction | UserContextMenuCommandInteraction | MessageContextMenuCommandInteraction): Promise<string> {
        let content: string;
        if (interaction.guild) {
            try {
                const member = await interaction.guild.members.fetch(interaction.options.getUser("user")!);
                content = member.displayAvatarURL({ extension: "png", size: 1024 });
            } catch {
                content = "This member doesn't seem to be here. If you want to petpet them, do it in my DMs.";
            }
        } else {
            content = interaction.options.getUser("user")!.avatarURL({ extension: "png", size: 1024 })!;
        }
        return content;
    }

    async function handleMessageContextMenu(interaction: MessageContextMenuCommandInteraction, client: Client) {
        await interaction.deferReply();
        let content: string;
        if (interaction.guild) {
            try {
                const member = await interaction.guild.members.fetch(interaction.targetMessage.author.id);
                content = member.displayAvatarURL({ extension: "png", size: 1024 });
            } catch {
                content = "This member doesn't seem to be here. If you want to petpet them, do it in my DMs.";
            }
        } else {
            const user = await client.users.fetch(interaction.targetMessage.author.id);
            content = user.avatarURL({ extension: "png", size: 1024 })!;
        }
        if (!content.startsWith("http")) {
            return interaction.editReply({
                content: content
            });
        }
        const gif = await getPetGif(content);
        if (typeof gif === "string") {
            return interaction.editReply({ content: "I fucked up" });
        }
        sendGif(interaction, gif)
    }

    async function handleSlashCommand(interaction: ChatInputCommandInteraction) {
        const ephemeral = interaction.options.getBoolean("ephemeral") || false;
        if (interaction.commandName === "invite") {
            invite(interaction);
            return;
        }
        if (interaction.commandName === "update-counter") {
            updateCounter(interaction);
            return;
        }
        if (interaction.commandName === "live-counter") {
            liveCounter(interaction);
            return;
        }
        await interaction.deferReply({ ephemeral: ephemeral });
        let content = await getSlashURL(interaction);
        if (!content.startsWith("http")) {
            return interaction.editReply({
                content: content
            })
        }
        let gif = await createGif(interaction, content);
        if (typeof gif === "boolean") return;
        sendGif(interaction, gif)
    }

    async function handleUserContextMenu(interaction: UserContextMenuCommandInteraction, client: Client) {
        await interaction.deferReply();
        let content: string;
        if (interaction.guild) {
            try {
                const member = await interaction.guild.members.fetch(interaction.targetId);
                content = member.displayAvatarURL({ extension: "png", size: 1024 });
            } catch {
                content = "This member doesn't seem to be here. If you want to petpet them, do it in my DMs.";
            }
        } else {
            const user = await client.users.fetch(interaction.targetId);
            content = user.avatarURL({ extension: "png", size: 1024 })!;
        }
        if (!content.startsWith("http")) {
            return interaction.editReply({
                content: content
            });
        }
        const gif = await getPetGif(content);
        if (typeof gif === "string") {
            return interaction.editReply({ content: "I fucked up" });
        }
        sendGif(interaction, gif)
    }

    function invite(interaction: ChatInputCommandInteraction) {
        return interaction.reply({
            content: "Click the button below to add me to one of your servers, or share this link to your friends to invite me: <https://lumap-is.gay/petpet>",
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
    }

    function updateCounter(interaction: ChatInputCommandInteraction) {
        if (interaction.user.id !== '635383782576357407') return interaction.reply({ content: "no", ephemeral: true });
        petCounter = interaction.options.getInteger("count")!;
        setActivity(client);
        return interaction.reply({
            content: "Alr, updated the counter to **" + petCounter + "**",
            ephemeral: true
        });
    }

    function liveCounter(interaction: ChatInputCommandInteraction) {
        return interaction.reply({
            content: "Since February 2023, the bot has been used **" + petCounter + "** times!"
        });
    }

    function isImage(url: string) {
        return /\.(jpg|jpeg|png)$/.test(url.toLowerCase());
    }

    function isRatelimited(interaction: CommandInteraction): boolean {
        const userRatelimits = rateLimits.filter(c => c.time + 60000 > Date.now()).filter(c => c.id === interaction.user.id);
        if (userRatelimits.length > 5) {
            interaction.reply({
                ephemeral: true,
                content: "You've been ratelimited (1 petpet per 10 seconds). This system has been put in place to prevent spam and power consuption from my server. Wait ~1 minute before being able to use the bot again."
            });
            return true;
        }
        return false;
    }

    function sendGif(interaction: CommandInteraction, gif: Buffer) {
        interaction.editReply({
            files: [
                {
                    attachment: gif,
                    name: "pet.gif",
                    description: "Pet!"
                }
            ],
        });
        addPetCounter()
    }

    function addPetCounter() {
        petCounter++;
    }


    client.on("interactionCreate", async (interaction) => {
        if (!interaction.isCommand()) return;
        if (isRatelimited(interaction)) return;
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
        client.user!.setPresence({ activities: [{ name: `${petCounter} petpets`, type: ActivityType.Watching }] });
        axios.post(config.petCounterWebhook, {
            content: `Last known petpet count: **${petCounter}**`
        })
    }
})();