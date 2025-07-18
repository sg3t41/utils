
# ----------------------------------------------------------------
# 1. Build Stage
# ----------------------------------------------------------------
# Gradle公式イメージを使用し、ビルド環境としてJDK 21を指定
FROM gradle:8.9-jdk21 AS build

# 作業ディレクトリを設定
WORKDIR /app

# まずは依存関係定義ファイルのみをコピー
# これにより、ソースコードの変更がない限り、依存関係のダウンロードはキャッシュされる
COPY build.gradle.kts ./
COPY gradle ./gradle

# 依存関係をダウンロード
RUN gradle dependencies --no-daemon

# アプリケーションのソースコードをコピー
COPY src ./src

# アプリケーションをビルドする（テストはスキップ）
RUN gradle build --no-daemon -x test

# ----------------------------------------------------------------
# 2. Package Stage
# ----------------------------------------------------------------
# 実行用に、より軽量なJREイメージを使用
FROM eclipse-temurin:21-jre-jammy

# 作業ディレクトリを設定
WORKDIR /app

# ビルドステージから、ビルドされたJARファイルのみをコピー
COPY --from=build /app/build/libs/*.jar app.jar

# アプリケーションの起動コマンド
ENTRYPOINT ["java", "-jar", "app.jar"]
