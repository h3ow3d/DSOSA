from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    model_config = SettingsConfigDict(env_file=".env", extra="ignore")

    database_url: str = "postgresql+psycopg://dsovs:dsovs@db:5432/dsovs"
    dsovs_api_url: str = (
        "https://owasp.org/www-project-devsecops-verification-standard/dist/dsovs.json"
    )
    frontend_origin: str = "http://localhost:5173"
    auto_sync_catalogue: bool = False


settings = Settings()
