package com.gakaki.demo.model;

import lombok.Builder;
import lombok.Data;
import lombok.experimental.SuperBuilder;

import java.util.ArrayList;
import java.util.List;

@Data
@Builder
public class SaltTigerBookItem {
    private String id;
    private String url;
    private String yearmonth;
    private String title;
    private String thumbnil;
    private String pubDate;
    private String officalUrl;
    private String officalPress;
    private String baiduUrl;
    private List<String> otherLinks;
    private String baiduCode;
    private String description;
    private String createdAt;
    @Builder.Default
    private List<SaltTigerBookTag> tags = new ArrayList<>();
    private String zlibSearchUrl;
}
