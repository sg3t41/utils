package com.github.sg3t41.utils.healthcheck.interfaces;

import static org.mockito.Mockito.when;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;
import static org.springframework.test.web.servlet.result.MockMvcResultHandlers.print;

import java.time.LocalDateTime;
import java.time.temporal.ChronoUnit;

import org.junit.jupiter.api.BeforeEach; // Add this import
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc; // Add this import
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.MvcResult;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.github.sg3t41.utils.healthcheck.usecase.HealthCheckDto;
import com.github.sg3t41.utils.healthcheck.usecase.HealthCheckUsecase;

@WebMvcTest(HealthCheckController.class)
@AutoConfigureMockMvc(addFilters = false)
public class HealthCheckControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @MockBean
    private HealthCheckUsecase healthCheckUsecase;

    @Autowired
    private ObjectMapper objectMapper;

    // 固定のLocalDateTime値を使用
    private static final LocalDateTime FIXED_CHECKED_AT = LocalDateTime.of(2025, 1, 1, 10, 0, 0);

    @BeforeEach
    void setUp() {
        objectMapper.registerModule(new JavaTimeModule());
        objectMapper.disable(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS);
    }

    @Test
    void checkHealth_returnsOkStatusAndHealthDto() throws Exception {
        // Mocking the use case behavior
        HealthCheckDto mockDto = HealthCheckDto.builder()
                .status("OK")
                .checkedAt(FIXED_CHECKED_AT)
                .details("データベース接続は正常です。")
                .build();
        when(healthCheckUsecase.execute()).thenReturn(mockDto);

        // Perform the GET request and assert the response
        mockMvc.perform(get("/api/health"))
                .andDo(print()) // Print request and response details
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.status").value("OK")) // Move assertions here
                .andExpect(jsonPath("$.details").value("データベース接続は正常です。")) // Move assertions here
                .andReturn(); // Capture the MvcResult
    }

    @Test
    void checkHealth_returnsServiceUnavailableStatusWhenNotOk() throws Exception {
        // Mocking the use case behavior for an error scenario
        HealthCheckDto mockDto = HealthCheckDto.builder()
                .status("ERROR")
                .checkedAt(FIXED_CHECKED_AT)
                .details("データベース接続に失敗しました。")
                .build();
        when(healthCheckUsecase.execute()).thenReturn(mockDto);

        // Perform the GET request and assert the response
        mockMvc.perform(get("/api/health"))
                .andDo(print()) // Print request and response details
                .andExpect(status().isServiceUnavailable())
                .andExpect(jsonPath("$.status").value("ERROR")) // Move assertions here
                .andExpect(jsonPath("$.details").value("データベース接続に失敗しました。")) // Move assertions here
                .andReturn(); // Capture the MvcResult
    }
}
