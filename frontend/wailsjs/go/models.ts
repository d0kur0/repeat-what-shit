export namespace main {
	
	export interface MacroAction {
	    keys: number[];
	    delay: number;
	}
	export interface Macro {
	    id: string;
	    name: string;
	    activation_keys: number[];
	    type: number;
	    actions: MacroAction[];
	}
	export interface AppData {
	    macros: Macro[];
	}
	

}

