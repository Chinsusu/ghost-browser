export namespace ai {
	
	export class ActiveHours {
	    start: string;
	    end: string;
	
	    static createFrom(source: any = {}) {
	        return new ActiveHours(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.start = source["start"];
	        this.end = source["end"];
	    }
	}
	export class ScheduledActivity {
	    timeStart: string;
	    timeEnd: string;
	    activity: string;
	    sites: string[];
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new ScheduledActivity(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.timeStart = source["timeStart"];
	        this.timeEnd = source["timeEnd"];
	        this.activity = source["activity"];
	        this.sites = source["sites"];
	        this.description = source["description"];
	    }
	}
	export class DaySchedule {
	    dayOfWeek: number;
	    activities: ScheduledActivity[];
	
	    static createFrom(source: any = {}) {
	        return new DaySchedule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dayOfWeek = source["dayOfWeek"];
	        this.activities = this.convertValues(source["activities"], ScheduledActivity);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MouseProfile {
	    movementSpeed: string;
	    precision: string;
	    scrollBehavior: string;
	    clickDelay: number;
	    jitter: number;
	
	    static createFrom(source: any = {}) {
	        return new MouseProfile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.movementSpeed = source["movementSpeed"];
	        this.precision = source["precision"];
	        this.scrollBehavior = source["scrollBehavior"];
	        this.clickDelay = source["clickDelay"];
	        this.jitter = source["jitter"];
	    }
	}
	export class TypingProfile {
	    averageWpm: number;
	    variance: number;
	    errorRate: number;
	    pauseBetween: number;
	    thinkingPauses: boolean;
	
	    static createFrom(source: any = {}) {
	        return new TypingProfile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.averageWpm = source["averageWpm"];
	        this.variance = source["variance"];
	        this.errorRate = source["errorRate"];
	        this.pauseBetween = source["pauseBetween"];
	        this.thinkingPauses = source["thinkingPauses"];
	    }
	}
	export class WritingStyle {
	    formality: string;
	    tone: string;
	    verbosity: string;
	    useEmojis: boolean;
	    useSlang: boolean;
	    commonPhrases: string[];
	
	    static createFrom(source: any = {}) {
	        return new WritingStyle(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.formality = source["formality"];
	        this.tone = source["tone"];
	        this.verbosity = source["verbosity"];
	        this.useEmojis = source["useEmojis"];
	        this.useSlang = source["useSlang"];
	        this.commonPhrases = source["commonPhrases"];
	    }
	}
	export class Personality {
	    id: string;
	    profileId: string;
	    name: string;
	    age: number;
	    gender: string;
	    occupation: string;
	    location: string;
	    bio: string;
	    interests: string[];
	    expertiseAreas: string[];
	    writingStyle: WritingStyle;
	    typingSpeed: TypingProfile;
	    mouseBehavior: MouseProfile;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Personality(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.profileId = source["profileId"];
	        this.name = source["name"];
	        this.age = source["age"];
	        this.gender = source["gender"];
	        this.occupation = source["occupation"];
	        this.location = source["location"];
	        this.bio = source["bio"];
	        this.interests = source["interests"];
	        this.expertiseAreas = source["expertiseAreas"];
	        this.writingStyle = this.convertValues(source["writingStyle"], WritingStyle);
	        this.typingSpeed = this.convertValues(source["typingSpeed"], TypingProfile);
	        this.mouseBehavior = this.convertValues(source["mouseBehavior"], MouseProfile);
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Schedule {
	    id: string;
	    profileId: string;
	    timezone: string;
	    activeHours: ActiveHours;
	    weeklySchedule: DaySchedule[];
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Schedule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.profileId = source["profileId"];
	        this.timezone = source["timezone"];
	        this.activeHours = this.convertValues(source["activeHours"], ActiveHours);
	        this.weeklySchedule = this.convertValues(source["weeklySchedule"], DaySchedule);
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	

}

export namespace fingerprint {
	
	export class AudioFP {
	    noise: number;
	    sampleRate: number;
	
	    static createFrom(source: any = {}) {
	        return new AudioFP(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.noise = source["noise"];
	        this.sampleRate = source["sampleRate"];
	    }
	}
	export class CanvasFP {
	    noise: number;
	
	    static createFrom(source: any = {}) {
	        return new CanvasFP(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.noise = source["noise"];
	    }
	}
	export class MimeType {
	    type: string;
	    suffixes: string;
	
	    static createFrom(source: any = {}) {
	        return new MimeType(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.suffixes = source["suffixes"];
	    }
	}
	export class Plugin {
	    name: string;
	    filename: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new Plugin(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.filename = source["filename"];
	        this.description = source["description"];
	    }
	}
	export class MiscFP {
	    plugins: Plugin[];
	    mimeTypes: MimeType[];
	    permissions: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new MiscFP(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.plugins = this.convertValues(source["plugins"], Plugin);
	        this.mimeTypes = this.convertValues(source["mimeTypes"], MimeType);
	        this.permissions = source["permissions"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TimezoneFP {
	    timezone: string;
	    timezoneOffset: number;
	    locale: string;
	
	    static createFrom(source: any = {}) {
	        return new TimezoneFP(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.timezone = source["timezone"];
	        this.timezoneOffset = source["timezoneOffset"];
	        this.locale = source["locale"];
	    }
	}
	export class NetworkFP {
	    webRTCPolicy: string;
	    publicIP?: string;
	    localIPs?: string[];
	    connectionType: string;
	    effectiveType: string;
	    downlink: number;
	    rtt: number;
	
	    static createFrom(source: any = {}) {
	        return new NetworkFP(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.webRTCPolicy = source["webRTCPolicy"];
	        this.publicIP = source["publicIP"];
	        this.localIPs = source["localIPs"];
	        this.connectionType = source["connectionType"];
	        this.effectiveType = source["effectiveType"];
	        this.downlink = source["downlink"];
	        this.rtt = source["rtt"];
	    }
	}
	export class HardwareFP {
	    batteryCharging?: boolean;
	    batteryLevel?: number;
	
	    static createFrom(source: any = {}) {
	        return new HardwareFP(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.batteryCharging = source["batteryCharging"];
	        this.batteryLevel = source["batteryLevel"];
	    }
	}
	export class FontFP {
	    installedFonts: string[];
	
	    static createFrom(source: any = {}) {
	        return new FontFP(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.installedFonts = source["installedFonts"];
	    }
	}
	export class WebGLFP {
	    vendor: string;
	    renderer: string;
	    unmaskedVendor: string;
	    unmaskedRenderer: string;
	    noise: number;
	
	    static createFrom(source: any = {}) {
	        return new WebGLFP(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.vendor = source["vendor"];
	        this.renderer = source["renderer"];
	        this.unmaskedVendor = source["unmaskedVendor"];
	        this.unmaskedRenderer = source["unmaskedRenderer"];
	        this.noise = source["noise"];
	    }
	}
	export class ScreenFP {
	    width: number;
	    height: number;
	    availWidth: number;
	    availHeight: number;
	    colorDepth: number;
	    pixelDepth: number;
	    pixelRatio: number;
	
	    static createFrom(source: any = {}) {
	        return new ScreenFP(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.width = source["width"];
	        this.height = source["height"];
	        this.availWidth = source["availWidth"];
	        this.availHeight = source["availHeight"];
	        this.colorDepth = source["colorDepth"];
	        this.pixelDepth = source["pixelDepth"];
	        this.pixelRatio = source["pixelRatio"];
	    }
	}
	export class NavigatorFP {
	    userAgent: string;
	    appVersion: string;
	    platform: string;
	    vendor: string;
	    language: string;
	    languages: string[];
	    hardwareConcurrency: number;
	    deviceMemory: number;
	    maxTouchPoints: number;
	    productSub: string;
	    doNotTrack: string;
	    cookieEnabled: boolean;
	    webdriver: boolean;
	
	    static createFrom(source: any = {}) {
	        return new NavigatorFP(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.userAgent = source["userAgent"];
	        this.appVersion = source["appVersion"];
	        this.platform = source["platform"];
	        this.vendor = source["vendor"];
	        this.language = source["language"];
	        this.languages = source["languages"];
	        this.hardwareConcurrency = source["hardwareConcurrency"];
	        this.deviceMemory = source["deviceMemory"];
	        this.maxTouchPoints = source["maxTouchPoints"];
	        this.productSub = source["productSub"];
	        this.doNotTrack = source["doNotTrack"];
	        this.cookieEnabled = source["cookieEnabled"];
	        this.webdriver = source["webdriver"];
	    }
	}
	export class Fingerprint {
	    navigator: NavigatorFP;
	    screen: ScreenFP;
	    webgl: WebGLFP;
	    canvas: CanvasFP;
	    audio: AudioFP;
	    fonts: FontFP;
	    hardware: HardwareFP;
	    network: NetworkFP;
	    timezone: TimezoneFP;
	    misc: MiscFP;
	
	    static createFrom(source: any = {}) {
	        return new Fingerprint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.navigator = this.convertValues(source["navigator"], NavigatorFP);
	        this.screen = this.convertValues(source["screen"], ScreenFP);
	        this.webgl = this.convertValues(source["webgl"], WebGLFP);
	        this.canvas = this.convertValues(source["canvas"], CanvasFP);
	        this.audio = this.convertValues(source["audio"], AudioFP);
	        this.fonts = this.convertValues(source["fonts"], FontFP);
	        this.hardware = this.convertValues(source["hardware"], HardwareFP);
	        this.network = this.convertValues(source["network"], NetworkFP);
	        this.timezone = this.convertValues(source["timezone"], TimezoneFP);
	        this.misc = this.convertValues(source["misc"], MiscFP);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	
	
	
	
	

}

export namespace profile {
	
	export class CreateOptions {
	    fingerprint?: fingerprint.Fingerprint;
	    proxyId?: string;
	    notes?: string;
	    tags?: string[];
	    os?: string;
	    browser?: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fingerprint = this.convertValues(source["fingerprint"], fingerprint.Fingerprint);
	        this.proxyId = source["proxyId"];
	        this.notes = source["notes"];
	        this.tags = source["tags"];
	        this.os = source["os"];
	        this.browser = source["browser"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Profile {
	    id: string;
	    name: string;
	    fingerprint?: fingerprint.Fingerprint;
	    proxyId?: string;
	    dataDir: string;
	    notes: string;
	    tags: string[];
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	    // Go type: time
	    lastUsedAt?: any;
	
	    static createFrom(source: any = {}) {
	        return new Profile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.fingerprint = this.convertValues(source["fingerprint"], fingerprint.Fingerprint);
	        this.proxyId = source["proxyId"];
	        this.dataDir = source["dataDir"];
	        this.notes = source["notes"];
	        this.tags = source["tags"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.lastUsedAt = this.convertValues(source["lastUsedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace proxy {
	
	export class CheckResult {
	    proxyId: string;
	    success: boolean;
	    latency: number;
	    ip?: string;
	    country?: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new CheckResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.proxyId = source["proxyId"];
	        this.success = source["success"];
	        this.latency = source["latency"];
	        this.ip = source["ip"];
	        this.country = source["country"];
	        this.error = source["error"];
	    }
	}
	export class Proxy {
	    id: string;
	    name: string;
	    type: string;
	    host: string;
	    port: number;
	    username?: string;
	    password?: string;
	    country?: string;
	    // Go type: time
	    lastCheckAt?: any;
	    lastCheckStatus: string;
	    lastCheckLatency: number;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Proxy(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.host = source["host"];
	        this.port = source["port"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.country = source["country"];
	        this.lastCheckAt = this.convertValues(source["lastCheckAt"], null);
	        this.lastCheckStatus = source["lastCheckStatus"];
	        this.lastCheckLatency = source["lastCheckLatency"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

