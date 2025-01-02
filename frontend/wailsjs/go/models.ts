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
	    include_titles: string;
	}
	export interface AppData {
	    macros: Macro[];
	}
	

}

