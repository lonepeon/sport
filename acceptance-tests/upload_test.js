const { above, attach, openBrowser, closeBrowser, goto, click, fileField, into, screenshot, setConfig, text, textBox, timeField, waitFor, write } = require("taiko");
const assert = require("assert").strict;
const settings = require("./settings").load();
const helpers = require("./helpers");
const path = require("path");

const browserArgs = {args:['--no-sandbox', '--disable-setuid-sandbox']};

describe("unauthenticated session", () => {
	beforeEach(async() => {
		await openBrowser(browserArgs);
		await goto(settings.url);
	});

	afterEach(async() => {
		await closeBrowser();
	});

	describe("upload activity", () => {
		helpers.uiTest("receive an error", async () => {
			await click("Upload Activity");
			await assert.ok(await text("Username:").exists());
		});
	})
});

describe("authenticated session", () => {
	beforeEach(async() => {
		await openBrowser(browserArgs);
		await goto(settings.url);
	});

	afterEach(async() => {
		await closeBrowser();
	});

	describe("upload activity", () => {
		helpers.uiTest("works", async () => {
			await click("Upload Activity");
			await write(settings.username, into(textBox("Username:")))
			await write(settings.password, into(textBox("Password:")))
			await click("Login");
			await timeField("Date:").select(new Date("2021-01-31T22:01:00"));
			await attach(path.join(__dirname, "testdata", "valid.gpx"), fileField("GPX file:"));
			await click("Submit");
			await helpers.predicateOrReload(async () => { return text("2021/01/31 22:01").exists() }, {retry: 10, timeout: 5000});
			await goto(settings.url + "/running-session/202101312201")
			await assert.ok(await text("4.18km").exists(0,0));
			await assert.ok(await text("10.03km/h").exists(0,0));
			await assert.ok(await text("25m0.607s").exists(0,0));
			await goto(settings.url);
			await click("Delete", above("2021/01/31 22:01"));
			await click("I confirm");
			await goto(settings.url);
			await helpers.predicateOrReload(async () => { return !await text("2021/01/31 22:01").exists()}, {retry: 10, timeout: 5000});
		});
	});
});
