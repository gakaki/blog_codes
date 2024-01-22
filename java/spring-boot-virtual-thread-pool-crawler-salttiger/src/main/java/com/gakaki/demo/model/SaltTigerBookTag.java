package com.gakaki.demo.model;

import lombok.Builder;
import lombok.Data;
import lombok.experimental.SuperBuilder;

@Data
@SuperBuilder
public class SaltTigerBookTag
{
    private String name;
    private String url;
}
