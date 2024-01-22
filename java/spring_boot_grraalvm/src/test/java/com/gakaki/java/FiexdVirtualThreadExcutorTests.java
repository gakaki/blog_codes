package com.gakaki.java;

import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;

class FixedVirtualThreadExecutorServiceTests {

	@Test
	void testExecutors() {

		ExecutorService executorService = new FixedVirtualThreadExecutorService(10);
		final RestTemplate restTemplate = new RestTemplate();
		for (int i=0;i<20;i++) {
			executorService.execute(() -> {
				System.out.println(restTemplate.getForObject("https://www.google.com?q=Shazin", String.class));
			});
		}

	}

}
