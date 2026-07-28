const fields = [
	"Passenger in wrong seat",
	"Crying baby",
	"Passanger in wrong boarding lane",
	"Queue at the toilets",
	"Bus to the airplane",
	"Delayed flight",
	"Personal item in overhead locker",
	"Boarding from wrong end",
	"Animal on flight",
	"Caught with overweight carry on",
	"Caught with oversized carry on",
	"Passenger smells sweaty",
	"Passenger buys small champagne",
	"Clapping for landing",
	"Passenger buys lottery ticket",
	"Crying baby on flight",
	"Passenger goes barefoot",
	"Passenger switches seats after takeoff",
	"Passenger takes a call on speaker",
	"Passenger uses the service call button",
	"Passenger reclines their seat",
	"Passenger kicks the back of the seat",
	"Passenger snores loudly",
	"Passenger spills a drink",
	"Passangers fight",
	"Flight is overbooked",
	"Full flight asked to check in carry on",
];

export function getFields(n: number) {
	if (n > fields.length) {
		throw new Error("n cant be larger then " + fields.length);
	}

	let i = 0;
	let output: string[] = Array(n);
	const used = new Set<number>();

	while (true) {
		const index = Math.floor(Math.random() * fields.length);
		if (!used.has(index)) {
			output[i] = fields[index];
			used.add(index);
			i++;
		}
		if (i === n) {
			break;
		}
	}

	return output;
}
