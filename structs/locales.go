package structs

type Locale struct {
	Locale       string `json:"locale"`
	LanguageName string `json:"language_name"`
	NativeName   string `json:"native_name"`
}

var (
	IndonesianLocale Locale = Locale{
		Locale:       "id",
		LanguageName: "Indonesian",
		NativeName:   "Bahasa Indonesia",
	}
	DanishLocale Locale = Locale{
		Locale:       "da",
		LanguageName: "Danish",
		NativeName:   "Dansk",
	}
	GermanLocale Locale = Locale{
		Locale:       "de",
		LanguageName: "German",
		NativeName:   "Deutsch",
	}
	EnglishUKLocale Locale = Locale{
		Locale:       "en-GB",
		LanguageName: "English, UK",
		NativeName:   "English, UK",
	}
	EnglishUSLocale Locale = Locale{
		Locale:       "en-US",
		LanguageName: "English, US",
		NativeName:   "English, US",
	}
	SpanishLocale Locale = Locale{
		Locale:       "es-ES",
		LanguageName: "Spanish",
		NativeName:   "Español",
	}
	SpanishLATAMLocale Locale = Locale{
		Locale:       "es-419",
		LanguageName: "Spanish, LATAM",
		NativeName:   "Español, LATAM",
	}
	FrenchLocale Locale = Locale{
		Locale:       "fr",
		LanguageName: "French",
		NativeName:   "Français",
	}
	CroatianLocale Locale = Locale{
		Locale:       "hr",
		LanguageName: "Croatian",
		NativeName:   "Hrvatski",
	}
	ItalianLocale Locale = Locale{
		Locale:       "it",
		LanguageName: "Italian",
		NativeName:   "Italiano",
	}
	LithuanianLocale Locale = Locale{
		Locale:       "lt",
		LanguageName: "Lithuanian",
		NativeName:   "Lietuviškai",
	}
	HungarianLocale Locale = Locale{
		Locale:       "hu",
		LanguageName: "Hungarian",
		NativeName:   "Magyar",
	}
	DutchLocale Locale = Locale{
		Locale:       "nl",
		LanguageName: "Dutch",
		NativeName:   "Nederlands",
	}
	NorwegianLocale Locale = Locale{
		Locale:       "no",
		LanguageName: "Norwegian",
		NativeName:   "Norsk",
	}
	PolishLocale Locale = Locale{
		Locale:       "pl",
		LanguageName: "Polish",
		NativeName:   "Polski",
	}
	PortugueseBrazilianLocale Locale = Locale{
		Locale:       "pt-BR",
		LanguageName: "Portuguese, Brazilian",
		NativeName:   "Português do Brasil",
	}
	RomanianLocale Locale = Locale{
		Locale:       "ro",
		LanguageName: "Romanian, Romania",
		NativeName:   "Română",
	}
	FinnishLocale Locale = Locale{
		Locale:       "fi",
		LanguageName: "Finnish",
		NativeName:   "Suomi",
	}
	SwedishLocale Locale = Locale{
		Locale:       "sv-SE",
		LanguageName: "Swedish",
		NativeName:   "Svenska",
	}
	VietnameseLocale Locale = Locale{
		Locale:       "vi",
		LanguageName: "Vietnamese",
		NativeName:   "Tiếng Việt",
	}
	TurkishLocale Locale = Locale{
		Locale:       "tr",
		LanguageName: "Turkish",
		NativeName:   "Türkçe",
	}
	CzechLocale Locale = Locale{
		Locale:       "cs",
		LanguageName: "Czech",
		NativeName:   "Čeština",
	}
	GreekLocale Locale = Locale{
		Locale:       "el",
		LanguageName: "Greek",
		NativeName:   "Ελληνικά",
	}
	BulgarianLocale Locale = Locale{
		Locale:       "bg",
		LanguageName: "Bulgarian",
		NativeName:   "български",
	}
	RussianLocale Locale = Locale{
		Locale:       "ru",
		LanguageName: "Russian",
		NativeName:   "Pусский",
	}
	UkrainianLocale Locale = Locale{
		Locale:       "uk",
		LanguageName: "Ukrainian",
		NativeName:   "Українська",
	}
	HindiLocale Locale = Locale{
		Locale:       "hi",
		LanguageName: "Hindi",
		NativeName:   "हिन्दी",
	}
	ThaiLocale Locale = Locale{
		Locale:       "th",
		LanguageName: "Thai",
		NativeName:   "ไทย",
	}
	ChineseChinaLocale Locale = Locale{
		Locale:       "zh-CN",
		LanguageName: "Chinese, China",
		NativeName:   "中文",
	}
	JapaneseLocale Locale = Locale{
		Locale:       "ja",
		LanguageName: "Japanese",
		NativeName:   "日本語",
	}
	ChineseTaiwanLocale Locale = Locale{
		Locale:       "zh-TW",
		LanguageName: "Chinese, Taiwan",
		NativeName:   "繁體中文",
	}
	KoreanLocale Locale = Locale{
		Locale:       "ko",
		LanguageName: "Korean",
		NativeName:   "한국어",
	}
)

// AllLocales is a slice containing all Locale variables
var AllLocales = []Locale{
    IndonesianLocale,
    DanishLocale,
    GermanLocale,
    EnglishUKLocale,
    EnglishUSLocale,
    SpanishLocale,
    SpanishLATAMLocale,
    FrenchLocale,
    CroatianLocale,
    ItalianLocale,
    LithuanianLocale,
    HungarianLocale,
    DutchLocale,
    NorwegianLocale,
    PolishLocale,
    PortugueseBrazilianLocale,
    RomanianLocale,
    FinnishLocale,
    SwedishLocale,
    VietnameseLocale,
    TurkishLocale,
    CzechLocale,
    GreekLocale,
    BulgarianLocale,
    RussianLocale,
    UkrainianLocale,
    HindiLocale,
    ThaiLocale,
    ChineseChinaLocale,
    JapaneseLocale,
    ChineseTaiwanLocale,
    KoreanLocale,
}

// FindLocaleByCode checks if a given string matches any of the Locale.Locale properties
func FindLocaleByCode(code string) *Locale {
    for _, locale := range AllLocales {
        if locale.Locale == code {
            return &locale
        }
    }
    return nil
}
