import numpy as np
import pandas as pd
import copy
import re
import math
from scipy import spatial
from sklearn.neighbors import NearestNeighbors
import grpc
import movieapi_pb2
import movieapi_pb2_grpc
from concurrent import futures

"""# Loading Dataset

A total of 8807 Movies/TV Shows
"""

netflix_df = pd.read_csv("./netflix_titles.csv")
netflix_df.shape

netflix_df.head()

# Replacing all NaN values with "missing"
netflix_df.fillna('missing', inplace = True)
netflix_df.head()

netflix_df.info()

# Checking if there are any NaN still left
netflix_df.isnull().sum()

"""# Data Preprocessing and Cleaning"""

# Changing the Column "Listed" to "Genre"
netflix_df.rename(columns= {"listed_in": "genre"}, inplace=True)
netflix_df.head()

"""Columns Considered for Predictions are as Follows:

1. Country
2. Release_Year
3. Rating
4. Duration
5. Genre 
"""

recommendation_cols = ["country", "release_year", "rating", "duration", "genre"]
new_df = copy.deepcopy(netflix_df[recommendation_cols])
new_df.head()

country = []
release_year = [] 
rating = []
duration = [] 
genres = []

def split_by_delimeters(target_list):
    """
    this method splits a target list by some delimeters
    """
    result = []
    for i in target_list:
        delimiters = ",", "&"
        regex_pattern = '|'.join(map(re.escape, delimiters))
        result.extend(re.split(regex_pattern, i))
    result = [i.strip() if i not in ['', 'missing'] else i for i in result]
    return result

# preparing all columns for the dataset
country = list(set(split_by_delimeters(new_df['country'])))
release_year = list(set(new_df['release_year']))
release_year = [str(i) for i in release_year]
ratings = list(set(new_df['rating']))
seasons_durations = ['1_season', '2_season', '3_season', '4_season','5+_season']
movies_durations = ['0_25_min', '26_50_min', '51_75_min', '76_100_min', 
                    '101_125_min', '126_150_min', '151+_min' ]
durations = seasons_durations + movies_durations
genres = list(set(split_by_delimeters(new_df['genre'])))

# combining all columns for the one hot encoded vector form
all_columns = country + release_year + ratings + durations + genres
all_columns.remove('missing')

# initializes a df with '0' values for the one-hot-encoded vector
ohe_df = pd.DataFrame(0, index = np.arange(len(new_df)), columns = all_columns)

def duration_adjustment(duration: str) -> str:
    try:
        dur_list = []
        if 'Season' in duration:
            temp_res = duration.split()
            no_of_seasons = int(temp_res[0])
            if no_of_seasons <5:
                return seasons_durations[no_of_seasons - 1]
            return seasons_durations[-1]

        else:
            temp_res = duration.split()
            runtime_mins = int(temp_res[0])
            if runtime_mins <= 150:
                index = math.ceil((runtime_mins/25) - 1.0)
                return movies_durations[index]
            return movies_durations[-1]
    except:
        return 'missing'

def return_columns(row):
    """
    recieves a df row and returns the respective columns/features
    that the item i.e. movie falls in
    """
    result_cols = []
    result_cols.extend(split_by_delimeters([row['country']]))
    result_cols.extend(split_by_delimeters([row['genre']]))
    result_cols.append(str(row['release_year']))
    result_cols.append(row['rating'])
    result_cols.append(duration_adjustment(str(row['duration'])))
    if 'missing' in result_cols:
        result_cols.remove('missing')
    return result_cols

# preparing the one hot encoded df of all items i.e. movies as vectors
for ind,row in new_df.iterrows():
    ohe_df.loc[ind, return_columns(row)] = 1

ohe_df.head()

def recommend_by_cosine(movie, top_items):
    """
    recommends top_similar movies based on cosine similarity
    """
    movie_index = netflix_df[netflix_df['title'] == movie].index[0]
    vector = ohe_df.iloc[movie_index]
    distance = []
    for ind, row in ohe_df.iterrows():
        distance.append(spatial.distance.cosine(vector, row))
    
    indexes = sorted(range(len(distance)), key=lambda i: distance[i])[:top_items + 1]
    
    return list(netflix_df.iloc[indexes]['title'])[1:]

def recommend_by_knn(movie, top_items):
    """
    recommends top_similar movies based on knn algorithm
    """
    movie_index = netflix_df[netflix_df['title'] == movie].index[0]
    vector = ohe_df.iloc[movie_index]
    knn = NearestNeighbors(n_neighbors= top_items + 1, algorithm='auto')
    knn.fit(ohe_df.values)
    indexes = list(knn.kneighbors([vector], top_items + 1, return_distance=False)[0])
    return list(netflix_df.iloc[indexes]['title'])[1:]


class Movie_RecommendationServicer(movieapi_pb2_grpc.Movie_RecommendationServicer):
    def Movie_Recommendation_by_cnn(self, request, context):
        result = [request.moviename]
        result = (recommend_by_cosine(request.moviename, request.no_of_recommendations))
        return movieapi_pb2.RecommendedMovies_Array(Recommended_Movies=result)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    movieapi_pb2_grpc.add_Movie_RecommendationServicer_to_server(Movie_RecommendationServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()


