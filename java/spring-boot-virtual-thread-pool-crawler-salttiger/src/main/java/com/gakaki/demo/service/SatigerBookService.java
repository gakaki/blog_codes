package com.gakaki.demo.service;

import cn.hutool.core.util.ObjectUtil;
import com.gakaki.demo.model.SaltTigerBookItem;
import com.gakaki.demo.model.SaltTigerBookTag;
import lombok.SneakyThrows;
import lombok.extern.slf4j.Slf4j;
import org.jsoup.Jsoup;
import org.jsoup.internal.StringUtil;
import org.jsoup.nodes.Document;
import org.jsoup.nodes.Element;
import org.jsoup.select.Elements;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.*;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import java.util.stream.Collectors;

@Slf4j
@Service
public class SatigerBookService {

    @SneakyThrows
    public List<SaltTigerBookItem> fetchBooks() {

        String url = "https://salttiger.com/archives/";
        Document doc = Jsoup.connect(url).get();

        List<SaltTigerBookItem> items = new ArrayList<>();
        Elements listItems = doc.select("ul.car-list li"); // Example selector

        for (Element listItem : listItems) {
            String createdAt = listItem.select("span.car-yearmonth").text();
            Elements bookLinks = listItem.select("ul.car-monthlisting li a");

            for (Element bookLink : bookLinks) {
                String title = bookLink.text();
                String bookUrl = bookLink.attr("href");
//                System.out.println("Title: " + title + ", URL: " + bookUrl + ", Created At: " + createdAt);

                SaltTigerBookItem item = SaltTigerBookItem.builder().build();
                item.setTitle(title);
                item.setUrl(bookUrl);
                items.add(item);
            }
        }
        return items;
    }

    public static String regexFind(String regexStr, String str) {
        Pattern pattern = Pattern.compile(regexStr);
        Matcher matcher = pattern.matcher(str);
        if (matcher.find()) {
            return matcher.group();
        }
        return null;
    }

    @SneakyThrows
    public SaltTigerBookItem fetchDetail(SaltTigerBookItem item) {
        try {
            Document doc = Jsoup.connect(item.getUrl()).get();
            Elements articles = doc.select("article");
            for (Element article : articles) {
                item.setId(article.id());
                item.setThumbnil(article.select("div > p:nth-child(1) > img").attr("src"));
                String tmpText = article.select("strong").text();
                item.setPubDate(regexFind("\\d{4}\\.\\d{1,2}", tmpText)); // 出版时间：2020.12
                Element officalA = article.select("div > p:nth-child(1) > strong > a:nth-child(2)").first();

                if (ObjectUtil.isNotNull(officalA)) {
                    item.setOfficalPress(officalA.text());
                    item.setOfficalUrl(officalA.attr("href"));
                }

                for (Element it : article.select("article strong > a[href*=ed2k]")) {
                    item.getOtherLinks().add(it.attr("href"));
                }
                Element officalBaidu = article.select("article strong > a[href*=baidu]").first();

                if (ObjectUtil.isNotNull(officalBaidu)) {
                    item.setBaiduUrl(officalBaidu.attr("href"));
                    if (item.getBaiduUrl().contains("pwd")) {
                        item.setBaiduCode(regexFind("pwd=.*", item.getBaiduUrl()).replace("pwd=", ""));
                    } else {
                        item.setBaiduCode(regexFind("提取码    ：\\w{1,4}", tmpText).replace("提取码    ：", ""));
                    }
                }

                var description = regexFind("<p>内容简介([\\s\\S]*)", article.select("div.entry-content").html());
                if (!StringUtil.isBlank(description)) {
                    item.setDescription(description.replace("<p>内容简介：</p>", ""));
                }

                item.setCreatedAt(article.select("footer > a:nth-child(1) > time").attr("datetime"));
                item.setZlibSearchUrl(String.format("https://zlibrary-asia.se/s/%s?", item.getTitle()));
                for (Element e : article.select("footer > a[rel*=tag]")) {
                    SaltTigerBookTag tag = SaltTigerBookTag.builder().build(); // 假设您有一个Tag类来存储标签数据
                    tag.setUrl(e.attr("href"));
                    tag.setName(e.text());
                    item.getTags().add(tag);
                }
                // 假设您有一个方法来处理JSON输出
//                JSONArray salttigerItems = new JSONArray();
//                JSONObject itemJson = item.toJson();
//                salttigerItems.put(itemJson);
//                // 写入JSON数据到文件
//                writeToJsonFile(salttigerItems.toString(), "saltTiger.json");
//                // 处理zlibrary链接
//                List<String> totalZlibraryLinks = new ArrayList<>();
//                for (Object salttigerItemObj : salttigerItems.toList()) {
//                    JSONObject salttigerItem = (JSONObject) salttigerItemObj;
//                    totalZlibraryLinks.add(salttigerItem.getString("zlibSearchUrl"));
//                }
//                writeToJsonFile(new JSONArray(totalZlibraryLinks).toString(), "zlibrary.json");
            }
        } catch (RuntimeException e) {
            e.printStackTrace();
        }
        return item;
    }

    @SneakyThrows
    public List<SaltTigerBookItem> fetchBookDetails() {

        List<SaltTigerBookItem> items = this.fetchBooks();
        ExecutorService executorService = new FixedVirtualThreadExecutorService(3);
        //为了 代码简单 只用30个
        items = items.stream().limit(15).collect(Collectors.toList());

        // 创建一个任务列表
//        List<Callable<SaltTigerBookItem>> tasks = items.stream().map(item ->
//                new Callable<SaltTigerBookItem>() {
//                    @Override
//                    public SaltTigerBookItem call() throws Exception {
//                        return fetchDetail(item);
//                    }
//        }).collect(Collectors.toList());

        List<Callable<SaltTigerBookItem>> tasks = items.stream()
                .map(item -> (Callable<SaltTigerBookItem>) () -> fetchDetail(item))
                .collect(Collectors.toList());

        List<SaltTigerBookItem> itemsFinal = new ArrayList<>();

        // 使用invokeAll执行所有任务，并等待它们完成
        List<Future<SaltTigerBookItem>> futures = executorService.invokeAll(tasks);

        // 遍历Future列表，获取每个任务的结果
        for (Future<SaltTigerBookItem> future : futures) {
            try {
                System.out.println(future.get());
                itemsFinal.add(future.get());
            } catch (InterruptedException | ExecutionException e) {
                // 异常处理逻辑
                e.printStackTrace();
            }

        }
        executorService.shutdown();
        return items;

    }
}
