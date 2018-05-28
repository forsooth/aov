import Html exposing (..)
import Html.Attributes exposing (class, href)
import Http
import List exposing (head)
import Html.Events exposing (onClick)
import Json.Decode exposing (Decoder, at, list, int, string, succeed, field, map3, map5)
import Json.Decode.Pipeline exposing (decode, required, optional, hardcoded)
import Bootstrap.CDN as CDN
import Bootstrap.Grid as Grid
import Bootstrap.Button as Button

-- Model
type alias Model = List Feed
type alias Items = List Item

type alias Feed = {id : Int,
           title : String,
           description : String,
           link : String,
           items : List Item
          }

type alias Item = {title : String, link : String, updated : String}

-- Update

update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
    case msg of
    NoOp ->
        (model, Cmd.none)
    GetFeed ->
        (model, Http.send GotFeed getFeed)
    GotFeed result ->
        case result of
        Err httpError ->
            let _ = Debug.log "handleError" httpError
            in (model, Cmd.none)
        Ok data ->
            (data, Cmd.none)

-- View

view : Model -> Html Msg
view model = Grid.container [class "text-center jumbotron"]
                [CDN.stylesheet, mainContent model]

-- API Things

type Msg
    = GotFeed (Result Http.Error Model)
    | GetFeed
    | NoOp

getFeedsCmd : Cmd Msg
getFeedsCmd = Http.send GotFeed getFeed

init : (Model, Cmd Msg)
init = ([], getFeedsCmd)

api : String
api = "http://localhost:8080/v1/feeds"

getFeed : Http.Request Model
getFeed = Http.get api feedListDecoder

-- Page Content
itemDecoder: Decoder Item
itemDecoder =
  map3 Item
    (field "title" string)
    (field "link" string)
    (field "updated" string)

feedDecoder: Decoder Feed
feedDecoder =
  map5 Feed
    (field "id" int)
    (field "title" string)
    (field "description" string)
    (field "link" string)
    (field "items" itemListDecoder)

itemListDecoder : Decoder Items
itemListDecoder =
  list itemDecoder

feedListDecoder : Decoder Model
feedListDecoder =
  list feedDecoder

feedHTML : Feed -> Html Msg
feedHTML feed = 
    div [] [h5 [] [text "Most recent posts from "],
            h3 [] [text feed.title],
            div [class "row"] (List.map itemHTML feed.items)]
        

itemHTML : Item -> Html Msg
itemHTML item = 
    div [class "col-sm-4"]
        [h4 [class "alert alert-info"]
            [a [href item.link] [text item.title]]]

mainContent : Model -> Html Msg
mainContent model =
    div [] [div [class "jumbotron"]
                [h1 [] [text "RSS Reader"],
                 Button.button [Button.outlinePrimary,
                 Button.attrs [onClick GetFeed]] [text "Update"]],
            div [class "container-fluid"] (List.map feedHTML model)]

-- Run

main : Program Never Model Msg
main = Html.program {init = init,
             update = update,
             view = view,
             subscriptions = always Sub.none}



