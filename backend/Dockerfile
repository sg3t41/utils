
# Build stage
FROM gradle:8.8.0-jdk21 AS build
WORKDIR /home/gradle/project
COPY . .
RUN gradle build --no-daemon

# Package stage
FROM eclipse-temurin:21-jre-jammy
WORKDIR /app
COPY --from=build /home/gradle/project/build/libs/*.jar app.jar
ENTRYPOINT ["java", "-jar", "app.jar"]
