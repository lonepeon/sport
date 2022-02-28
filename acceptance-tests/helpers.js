const { reload, screenshot, text } = require('taiko');
const assert = require("assert").strict;

exports.uiTest = async (name, fn) => {
	test(name, async () => {
		try {
			await fn();
		} catch (err) {
			await screenshot();
			throw err
		}
	});
};

exports.predicateOrReload = async (predicate, options) => {
	const sleep = async(ms) => {
		return new Promise(resolve => setTimeout(resolve, ms));
	}

	let i = 0;
	while (i < options.retry) {
		if (await predicate()) {
			return;
		} else {
			await sleep(options.timeout);
			await reload();
			i++;
		}
	}

	try {
		assert.equal('', txt)
	} catch(err) {
		Error.captureStackTrace(err, predicateOrReload)
		throw err
	}
}
