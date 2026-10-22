declare global {
	namespace App {
		// interface Error {}
		// interface Locals {}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
	
	type ResrapTemplate = {
		name: string;
		language: string;
		description: string;
		grammar: string;
	};

	var generateText:
		| ((grammar: string, startpoint: string, seed: number, tokenlen: number) => string)
		| undefined;
	var generateTemplate:
		| ((template: string, startpoint: string, seed: number, tokenlen: number) => string)
		| undefined;
	var getTemplates: (() => string) | undefined;
	var Go: new () => any;
	var goTest: () => string;
}

export {};
