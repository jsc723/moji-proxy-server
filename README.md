# moji proxy server
## build
```
go build -o server.exe *.go
```
## run
```
./server
./server --username <moji-username> --password <moji-password>
```
## api
### /search
POST
#### request
```json
{
    "query": "<the word to search>"
}
```
#### result
```
{
    "code": 200,
    "result": {
        "grammar": {
            "searchResult": [
                {
                    "excerpt": "[文法・文法・文法] 前项刚一实现，就马上进行后项，后项不能用过去时。\n\n接续：\n动词ます形+次第 ",
                    "libId": null,
                    "targetId": "LrjnnuN8bV",
                    "targetType": 106,
                    "title": "～次第 | ～しだい"
                },
                {
                    "excerpt": "[文法] ～が変われば結果も変わる・～によって決まる。\n如果……变化，那么结果也发生变化，根据……而",
                    "targetId": "CAmBGDyZSW",
                    "targetType": 106,
                    "title": "～次第だ | ～しだいだ"
                }
            ]
        },
        "word": {
            "searchResult": [
                {
                    "excerpt": "[名詞・接続詞] 顺序。（順序。） 情况，缘由。（由来。経過。なりゆき。いきさつ。） 听其自然，听任",
                    "isFree": true,
                    "targetId": "198961275",
                    "targetType": 102,
                    "title": "次第 | しだい ◎"
                },
                {
                    "excerpt": "[副] 逐渐，渐渐，慢慢。（事態が時の経過とともに少しずつ変化するさま。） ",
                    "isFree": true,
                    "targetId": "198961277",
                    "targetType": 102,
                    "title": "次第に | しだいに ◎"
                },
                {
                    "excerpt": "[文法] ～が実現した後、すぐに続けてある行動をする。硬い言い方。\n……实现以后，马上采取后续行为。",
                    "targetId": "6vHbBjPw2D",
                    "targetType": 102,
                    "title": "～次第 | ～しだい"
                },
                {
                    "excerpt": "[文法] ～が変われば結果も変わる・～によって決まる。\n如果……变化，那么结果也发生变化，根据……而",
                    "targetId": "CAmBGDyZSW",
                    "targetType": 102,
                    "title": "～次第だ | ～しだいだ"
                },
                {
                    "excerpt": "[文法・文法・文法] 前项刚一实现，就马上进行后项，后项不能用过去时。\n\n接续：\n动词ます形+次第 ",
                    "libId": null,
                    "targetId": "LrjnnuN8bV",
                    "targetType": 102,
                    "title": "～次第 | ～しだい"
                }
            ]
        }
    }
}
```
### /details
POST
#### request
```json
{
    "objectIds": [
        "<id returned by /search>",
        "<id returned by /search>",
        ...
    ]
}
```
#### result
```json
{
    "words": [
        {
            "id": "198970381",
            "spell": "食べる",
            "pron": "たべる",
            "accent": "②",
            "excerpt": "[他动#二类・惯用语] 吃。（飲食物をいただく。） 生活。（生計を立てる。） 惯用句。",
            "details": [
                {
                    "id": "58672",
                    "title": "他动#二类"
                },
                {
                    "id": "58673",
                    "title": "惯用语"
                }
            ],
            "subDetails": [
                {
                    "id": "81748",
                    "title": "惯用句。",
                    "detailId": "58673",
                    "examples": []
                },
                {
                    "id": "81746",
                    "title": "吃。（飲食物をいただく。）",
                    "detailId": "58672",
                    "examples": [
                        {
                            "id": "59443",
                            "title": "ご飯を食べる。",
                            "trans": "吃饭。"
                        },
                        {
                            "id": "59444",
                            "title": "腹いっぱい食べた。",
                            "trans": "吃饱了。"
                        },
                        {
                            "id": "59445",
                            "title": "わたしは一日中なにも食べなかった。",
                            "trans": "我整天什么也没有吃。"
                        },
                        {
                            "id": "59446",
                            "title": "くちゃくちゃ食べる。",
                            "trans": "咕唧咕唧地吃。"
                        },
                        {
                            "id": "59447",
                            "title": "ぼそぼそ食べる。",
                            "trans": "干巴巴地吃。"
                        },
                        {
                            "id": "59448",
                            "title": "ちびちび食べる。",
                            "trans": "一点一点地吃。"
                        },
                        {
                            "id": "14NK75eP2L",
                            "title": "食べてすぐ寝ると牛になる。",
                            "trans": "吃完饭马上躺下会变成牛。用于告诫不能做没有礼仪的事。"
                        }
                    ]
                },
                {
                    "id": "81747",
                    "title": "生活。（生計を立てる。）",
                    "detailId": "58672",
                    "examples": [
                        {
                            "id": "59449",
                            "title": "月給で食べる。",
                            "trans": "靠工资维持生活。"
                        },
                        {
                            "id": "59450",
                            "title": "この収入では食べられない。",
                            "trans": "靠这点收入生活维持不了。"
                        },
                        {
                            "id": "59451",
                            "title": "家族を食べさせていくために彼女は小さな店をやっていた。",
                            "trans": "她为了维持家庭生活，开了个小铺子。"
                        }
                    ]
                }
            ]
        },
        {
            "id": "L8Bwll55wV",
            "spell": "〜によって",
            "pron": "〜によって",
            "accent": "",
            "excerpt": "[文法] 通过、靠\n\n名詞に＋よって",
            "details": [
                {
                    "id": "vmXf1SaDEP",
                    "title": "文法"
                }
            ],
            "subDetails": [
                {
                    "id": "cKDyHbIROD",
                    "title": "通过、靠\n\n名詞に＋よって",
                    "detailId": "vmXf1SaDEP",
                    "examples": [
                        {
                            "id": "KeIQiZACGR",
                            "title": "いつの時代でも若者によって、新しい流行が作り出される。",
                            "trans": "不论哪个时代都靠年轻人来创造出新的潮流。"
                        },
                        {
                            "id": "jQfsiS8J4a",
                            "title": "関係者の皆様のご協力によって、無事この会を終了することができました。",
                            "trans": "通过所有相关人士的鼎力相助，会议才能顺利的结束。"
                        },
                        {
                            "id": "H8wIlrVxFj",
                            "title": "いくつかの国を旅してみて、食事の習慣が国によって違うことに驚いた。",
                            "trans": "去若干个国家旅行后，对饮食习惯因国而异感到惊讶。"
                        }
                    ]
                }
            ]
        },
        {
            "id": "LrjnnuN8bV",
            "spell": "～次第",
            "pron": "～しだい",
            "accent": "",
            "excerpt": "[文法・文法・文法] 前项刚一实现，就马上进行后项，后项不能用过去时。\n\n接续：\n动词ます形+次第 ",
            "details": [
                {
                    "id": "lszLGPUE4G",
                    "title": "文法"
                },
                {
                    "id": "fkWwEol6iq",
                    "title": "文法"
                },
                {
                    "id": "9orKTIeqGE",
                    "title": "文法"
                }
            ],
            "subDetails": [
                {
                    "id": "odYkdRCmJv",
                    "title": "前项刚一实现，就马上进行后项，后项不能用过去时。\n\n接续：\n动词ます形+次第",
                    "detailId": "lszLGPUE4G",
                    "examples": [
                        {
                            "id": "h2FHp65ixr",
                            "title": "検査の結果がわかり次第、ご連絡いたします。",
                            "trans": "检查的结果一出来，我们会马上与您联系。"
                        },
                        {
                            "id": "0V3cyWeDBk",
                            "title": "参加者の名前がわかり次第、教えていただけませんか。",
                            "trans": "知道参加者名字之后，您能不能马上告诉我？"
                        }
                    ]
                },
                {
                    "id": "2TUxIiTS4d",
                    "title": "表示情况，书面语。\n\n接续：\n动词/形容词/形容动词的名词修饰形+次第だ／次第で",
                    "detailId": "fkWwEol6iq",
                    "examples": [
                        {
                            "id": "I23egzpOXp",
                            "title": "先日お伝えした日程に誤りがありましたので、今回、改めてご連絡を差し上げた次第です。",
                            "trans": "由于前几天与您说的日期有误，所以现在再次和您联系。"
                        },
                        {
                            "id": "1C7FlvDitG",
                            "title": "以上のような次第で、来週の社員旅行は延期にさせていただきます。",
                            "trans": "由于上述情况，将下周的员工旅行延期。"
                        }
                    ]
                },
                {
                    "id": "3JOAUfUpqo",
                    "title": "表示根据，取决于。\n\n接续：\n名词+次第だ／次第で(は)",
                    "detailId": "9orKTIeqGE",
                    "examples": [
                        {
                            "id": "CFPnLguHv0",
                            "title": "あしたハイキングに行くかどうかは、お天気次第だ。",
                            "trans": "明天去不去郊游取决于天气。"
                        },
                        {
                            "id": "zeH6OzbaPe",
                            "title": "この製品は、アイデア次第でいろいろな使い方ができます。",
                            "trans": "这个产品，根据不同的构想可以有不同的使用方法。"
                        }
                    ]
                }
            ]
        }
    ]
}
```