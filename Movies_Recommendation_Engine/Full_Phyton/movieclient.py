from flask import Flask, request, render_template
import grpc
import movieapi_pb2
import movieapi_pb2_grpc

app = Flask(__name__)

@app.route('/')
def home():
    return render_template('index.html')

@app.route('/recommendations', methods=['POST'])
def get_recommendations():
    title = request.form['title']
    num_recommendations = request.form['num_recommendations']
    print(title)
    print(num_recommendations)
    recommended_movies = movieapi(title, num_recommendations)
    return render_template('recommendations.html', recommended_movies=recommended_movies)

def movieapi(st, n):
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = movieapi_pb2_grpc.Movie_RecommendationStub(channel)
        request = movieapi_pb2.MovieRecommendation(moviename=st, no_of_recommendations=n)
        response = stub.Movie_Recommendation_by_cnn(request)
        return response.Recommended_Movies

if __name__ == '__main__':
    print("HI")
    #rec = get_recommendations()
    #print(rec)
    app.run()
