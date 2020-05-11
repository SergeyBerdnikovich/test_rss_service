# test_rss_service
a test service for parsing rss for https://www.emerchantpay.com/

<h2>Basic Usage:</h2>
Run <code>make run</code>

<h2>Endpoints:</h2>
1) GET <code>/rss_feeds_items</code>
</br>
where url params: <code>array of urls with key "urls"</code>
</br>
header is jwt token: <code>Token: ...</code>
</br>
and response is JSON with array of items and errors messages if they are axists
</br>
<h4>Url example:</h4>
<code>/rss_feeds_items?urls=http://feeds.twit.tv/twit.xml&urls=http://rss.cnn.com/rss/cnn_topstories.rss</code>
<h4>Response structure example:</h4>
<pre>
{
    items: [
        {
            "title" => "Some item title",
            "source" => "Some feed title",
            "source_url" => "Some feed url",
            "link" => "Some item link",
            "publish_date" => "2020/05/01 22:45" // format can be different and the value can be nil
            "description" => "Some item description"
        },
        ...
    ],
    errors: "Error message 1, error message 2, ..."
}
</pre>
2) POST <code>/authenticate</code>
</br>
where body params: <code>login</code> and <code>password</code>
</br>
and response is JSON with login and jwt token or error message in case of wrong credentials
</br>
<h4>Response structure example:</h4>
<pre>
{
    "login" => "TestRssApp",
    "token" => "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODkyNzkwMDksImxvZ2luIjoiVGVzdFJzc0FwcCJ9.Qo2e_7IdxBBJVG76yEoH91cBdoIMErA4gFTp4bu2Hw4"
}
</pre>
